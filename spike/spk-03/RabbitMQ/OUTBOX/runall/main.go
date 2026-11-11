// SC-E runner: runs ALL scenarios (E1 to E6, B1, B2) one after another with one command and
// writes every log in one go: OUTBOX/<scenario>/{publisher,forwarder,consumer}_log.txt
// (the programs write their own logs, see the -log flag) and OUTBOX/runall_log.txt (the result table).
//
//	go run ./spike/spk-03/RabbitMQ/OUTBOX/runall
//
// Needs RabbitMQ and PostgreSQL reachable on localhost and the credentials in .env.
// For every scenario it purges the queue, starts the consumer, starts the forwarder when the
// scenario needs it, runs the publisher in the scenario's mode, waits until the consumer has
// been idle for a while, and checks the consumer's SUMMARY line against the expected values.
package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	rabbit "github.com/rabbitmq/amqp091-go"

	"github.com/ITA-Internship/GoJudge/spike/spk-03/RabbitMQ/OUTBOX/outlog"
	"github.com/ITA-Internship/GoJudge/spike/spk-03/RabbitMQ/rmqenv"
)

const (
	queue = "outbox_submissions"
	idle  = 10 * time.Second // the consumer stops after this long without a message
)

type scenario struct {
	name      string
	mode      string // publisher -mode
	forwarder string // "before" the publisher, "after" it, "none", or "crash" (killed after the first publish, then restarted)
	want      map[string]int
}

func want(deliveries, unique, dup, dbRows, dwm, mwd int) map[string]int {
	return map[string]int{"deliveries": deliveries, "unique": unique, "duplicates": dup,
		"db_rows": dbRows, "data_without_message": dwm, "message_without_data": mwd}
}

var scenarios = []scenario{
	{"E1-commit", "outbox", "before", want(5, 5, 0, 5, 0, 0)},
	{"E2-forwarder-down", "outbox", "after", want(5, 5, 0, 5, 0, 0)},
	{"E3-rollback", "outbox-rollback", "before", want(4, 4, 0, 4, 0, 0)},
	{"E4-crash-before-commit", "outbox-crash-before-commit", "before", want(2, 2, 0, 2, 0, 0)},
	{"E5-crash-after-commit", "outbox-crash-after-commit", "before", want(3, 3, 0, 3, 0, 0)},
	{"E6-forwarder-crash", "outbox", "crash", want(6, 5, 1, 5, 0, 0)},
	{"B1-direct-commit-first", "direct-commit-first", "none", want(2, 2, 0, 3, 1, 0)},
	{"B2-direct-publish-first", "direct-publish-first", "none", want(5, 5, 0, 4, 0, 1)},
}

var (
	out       io.Writer = os.Stdout
	outboxDir string
	rootDir   string
	binDir    string
)

func main() {
	_, thisFile, _, _ := runtime.Caller(0)
	outboxDir = filepath.Dir(filepath.Dir(thisFile))
	rootDir = filepath.Join(outboxDir, "..", "..", "..", "..")

	var closeLog func()
	out, closeLog = outlog.Open(filepath.Join(outboxDir, "runall_log.txt"))
	defer closeLog()

	var err error
	binDir, err = os.MkdirTemp("", "spk03-outbox")
	must(err)
	defer os.RemoveAll(binDir)
	for _, p := range []string{"consumer", "forwarder", "publisher"} {
		fmt.Fprintf(out, "building %s\n", p)
		cmd := exec.Command("go", "build", "-o", bin(p), "./spike/spk-03/RabbitMQ/OUTBOX/"+p)
		cmd.Dir = rootDir
		cmd.Stderr = os.Stderr
		must(cmd.Run())
	}

	// Rows left in the outbox by earlier runs would show up as extra messages: drain them first.
	fmt.Fprintln(out, "\n== prepare: drain the outbox, purge the queue")
	purge()
	f := startProc("forwarder", filepath.Join(outboxDir, "_prepare", "forwarder_log.txt"))
	waitLog(f.logPath, "forwarder started", 30*time.Second)
	time.Sleep(6 * time.Second)
	f.kill()
	purge()

	failed := 0
	var table []string
	for _, sc := range scenarios {
		got := runScenario(sc)
		ok := got != nil
		for k, v := range sc.want {
			if got[k] != v {
				ok = false
			}
		}
		verdict := "PASS"
		if !ok {
			verdict = "FAIL"
			failed++
		}
		table = append(table, fmt.Sprintf("%s  %-24s want=%s", verdict, sc.name, format(sc.want)))
		table = append(table, fmt.Sprintf("      %-24s got =%s", "", format(got)))
	}

	fmt.Fprintln(out, "\n================ RESULT ================")
	for _, l := range table {
		fmt.Fprintln(out, l)
	}
	if failed == 0 {
		fmt.Fprintln(out, "ALL SCENARIOS AS EXPECTED")
	} else {
		fmt.Fprintf(out, "%d SCENARIO(S) DIFFER FROM THE EXPECTED VALUES, see the logs in %s\n", failed, outboxDir)
	}
}

func runScenario(sc scenario) map[string]int {
	fmt.Fprintf(out, "\n== %s (publisher -mode=%s, forwarder: %s)\n", sc.name, sc.mode, sc.forwarder)
	dir := filepath.Join(outboxDir, sc.name)
	log := func(name string) string { return filepath.Join(dir, name) }
	purge()

	cons := startProc("consumer", log("consumer_log.txt"), "-idle", idle.String())
	waitLog(cons.logPath, "consumer started", 30*time.Second)

	var fwd, fwd2 *proc
	switch sc.forwarder {
	case "before":
		fwd = startProc("forwarder", log("forwarder_log.txt"))
		waitLog(fwd.logPath, "forwarder started", 30*time.Second)
	case "crash":
		fwd = startProc("forwarder", log("forwarder_log.txt"), "-crash-after-publish")
		waitLog(fwd.logPath, "forwarder started", 30*time.Second)
	}

	pub := startProc("publisher", log("publisher_log.txt"), "-mode", sc.mode)
	fmt.Fprintf(out, "   publisher finished: %v (an error is expected in the crash modes)\n", pub.wait(60*time.Second))

	switch sc.forwarder {
	case "after":
		fwd = startProc("forwarder", log("forwarder_log.txt"))
	case "crash":
		fmt.Fprintf(out, "   forwarder finished: %v (an error is expected, it kills itself)\n", fwd.wait(30*time.Second))
		fwd2 = startProc("forwarder", log("forwarder_restart_log.txt"))
	}

	if err := cons.wait(idle + 60*time.Second); err != nil {
		fmt.Fprintf(out, "   consumer: %v\n", err)
	}
	fwd.kill()
	fwd2.kill()

	got := summary(cons.logPath)
	fmt.Fprintf(out, "   consumer SUMMARY: %s\n", format(got))
	return got
}

// ---- helpers ----

type proc struct {
	cmd     *exec.Cmd
	logPath string
	done    chan error
}

// startProc starts one of the built programs. The program writes its own log file (-log).
func startProc(name, logPath string, args ...string) *proc {
	cmd := exec.Command(bin(name), append([]string{"-log", logPath}, args...)...)
	cmd.Dir = rootDir
	cmd.Stderr = os.Stderr // panics, if any
	must(cmd.Start())
	p := &proc{cmd: cmd, logPath: logPath, done: make(chan error, 1)}
	go func() { p.done <- cmd.Wait() }()
	return p
}

func (p *proc) wait(timeout time.Duration) error {
	select {
	case err := <-p.done:
		p.done <- err // allow waiting again
		return err
	case <-time.After(timeout):
		p.kill()
		return fmt.Errorf("timeout after %s, killed", timeout)
	}
}

func (p *proc) kill() {
	if p != nil && p.cmd.Process != nil {
		_ = p.cmd.Process.Kill()
	}
}

func waitLog(path, substr string, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if b, err := os.ReadFile(path); err == nil && strings.Contains(string(b), substr) {
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	panic(fmt.Sprintf("%q did not appear in %s within %s", substr, path, timeout))
}

// purge empties the destination queue (it is declared first, like the Watermill subscriber does).
func purge() {
	uri, err := rmqenv.URI()
	must(err)
	conn, err := rabbit.Dial(uri)
	must(err)
	defer conn.Close()
	ch, err := conn.Channel()
	must(err)
	defer ch.Close()
	_, err = ch.QueueDeclare(queue, true, false, false, false, nil)
	must(err)
	_, err = ch.QueuePurge(queue, false)
	must(err)
}

var kv = regexp.MustCompile(`(\w+)=(\d+)`)

// summary reads the consumer's SUMMARY line; nil if there is none.
func summary(path string) map[string]int {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	for _, line := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(line, "SUMMARY ") {
			m := map[string]int{}
			for _, x := range kv.FindAllStringSubmatch(line, -1) {
				m[x[1]], _ = strconv.Atoi(x[2])
			}
			return m
		}
	}
	return nil
}

func format(m map[string]int) string {
	if m == nil {
		return " (no SUMMARY line)"
	}
	return fmt.Sprintf(" deliveries=%d unique=%d duplicates=%d db_rows=%d data_without_message=%d message_without_data=%d",
		m["deliveries"], m["unique"], m["duplicates"], m["db_rows"], m["data_without_message"], m["message_without_data"])
}

func bin(name string) string {
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return filepath.Join(binDir, name)
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(out, "ERROR:", err)
		panic(err)
	}
}
