// Command outcap runs a go test -json (test2json) command and forwards its
// event stream to stdout with byte limits. Package-level events (start and
// the final pass/fail) are always forwarded, so the verdict survives an
// output flood. Test output beyond -limit and test-level events beyond
// -event-limit are dropped, and a single truncation marker is emitted.
package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// maxLine is the longest single event line we accept. Longer lines are
// dropped as a whole, so one giant line cannot exhaust memory.
const maxLine = 64 << 10

type event struct {
	Action string `json:"Action"`
	Test   string `json:"Test,omitempty"`
}

type marker struct {
	Action     string `json:"Action"`
	Output     string `json:"Output"`
	OutputType string `json:"OutputType"`
}

type limiter struct {
	out               *bufio.Writer
	outLimit, evLimit int64
	outBytes, evBytes int64
	dropped           int64
	truncated         bool
}

func (l *limiter) emit(text string) {
	b, _ := json.Marshal(marker{Action: "output", Output: text, OutputType: "error"})
	l.out.Write(append(b, '\n'))
}

func (l *limiter) drop(n int64) {
	l.dropped += n
	if !l.truncated {
		l.truncated = true
		l.emit(fmt.Sprintf("outcap: output limit reached, further test output is dropped (limit %d bytes)\n", l.outLimit))
	}
}

func (l *limiter) handle(line []byte) {
	n := int64(len(line))
	var ev event
	isOutput := json.Unmarshal(line, &ev) != nil || ev.Action == "output"

	switch {
	case !isOutput && ev.Test == "": // package-level start / pass / fail: always keep
		l.out.Write(line)
	case !isOutput: // test-level run / pass / fail / skip
		if l.evBytes+n <= l.evLimit {
			l.evBytes += n
			l.out.Write(line)
		} else {
			l.drop(n)
		}
	default: // test output, including frames and non-JSON noise
		if l.outBytes+n <= l.outLimit {
			l.outBytes += n
			l.out.Write(line)
		} else {
			l.drop(n)
		}
	}
}

func main() {
	outLimit := flag.Int64("limit", 1<<20, "max bytes of test output to forward")
	evLimit := flag.Int64("event-limit", 64<<10, "max bytes of test-level events to forward")
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "usage: outcap [-limit N] [-event-limit N] -- command [args...]")
		os.Exit(2)
	}

	// One pipe for both stdout and stderr of the child: nothing it writes
	// can reach the container log without passing through the limiter.
	pr, pw, err := os.Pipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "outcap: pipe: %v\n", err)
		os.Exit(1)
	}
	cmd := exec.Command(flag.Arg(0), flag.Args()[1:]...)
	cmd.Stdout, cmd.Stderr = pw, pw
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "outcap: start: %v\n", err)
		os.Exit(1)
	}
	pw.Close() // the child keeps its own copy; we only read

	// As PID 1 we must pass termination signals on to the child.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for s := range sigs {
			_ = cmd.Process.Signal(s)
		}
	}()

	l := &limiter{out: bufio.NewWriter(os.Stdout), outLimit: *outLimit, evLimit: *evLimit}
	r := bufio.NewReaderSize(pr, maxLine)
	for {
		line, err := r.ReadSlice('\n')
		if errors.Is(err, bufio.ErrBufferFull) {
			n := int64(len(line)) // an over-long line: skip it to its end
			for errors.Is(err, bufio.ErrBufferFull) {
				line, err = r.ReadSlice('\n')
				n += int64(len(line))
			}
			l.drop(n)
		} else if len(line) > 0 {
			l.handle(line)
		}
		if err != nil {
			if !errors.Is(err, io.EOF) {
				fmt.Fprintf(os.Stderr, "outcap: read: %v\n", err)
			}
			break
		}
	}

	waitErr := cmd.Wait()
	l.emit(fmt.Sprintf("outcap: forwarded %d bytes of output and %d bytes of test events, dropped %d bytes\n", l.outBytes, l.evBytes, l.dropped))
	l.out.Flush()

	code := 0
	var exitErr *exec.ExitError
	if errors.As(waitErr, &exitErr) {
		code = exitErr.ExitCode()
		if code < 0 { // killed by a signal
			code = 1
		}
	} else if waitErr != nil {
		code = 1
	}
	os.Exit(code)
}
