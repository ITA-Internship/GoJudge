//go:build linux

package trusted

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

const probeTarget = "1.1.1.1:80"

// httpHead sends one HTTP HEAD request over an existing connection
// and returns the response status line.
func httpHead(conn net.Conn) (string, error) {
	_ = conn.SetDeadline(time.Now().Add(3 * time.Second))
	if _, err := fmt.Fprint(conn, "HEAD / HTTP/1.1\r\nHost: 1.1.1.1\r\nConnection: keep-alive\r\n\r\n"); err != nil {
		return "", err
	}
	r := bufio.NewReader(conn)
	status, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	for { // skip headers until the empty line
		line, err := r.ReadString('\n')
		if err != nil {
			return "", err
		}
		if line == "\r\n" {
			break
		}
	}
	return strings.TrimSpace(status), nil
}

func TestEstablishedConnection(t *testing.T) {
	start := time.Now()
	at := func() time.Duration { return time.Since(start).Round(time.Millisecond) }

	// 1. Try to open a connection immediately after start
	conn, err := net.DialTimeout("tcp", probeTarget, 2*time.Second)
	if err != nil {
		t.Logf("t=%v: connection at start blocked: %v", at(), err)
		return // nothing was established, so nothing can be kept
	}
	defer conn.Close()
	t.Errorf("t=%v: connection at start succeeded (remote %s)", at(), conn.RemoteAddr())
	if status, err := httpHead(conn); err == nil {
		t.Logf("t=%v: request 1 over it: %s", at(), status)
	} else {
		t.Logf("t=%v: request 1 failed: %v", at(), err)
	}

	// 2. Wait until the policy rules have certainly been applied
	time.Sleep(5 * time.Second)

	// 3. A new connection must be blocked by now
	if c2, err := net.DialTimeout("tcp", probeTarget, 2*time.Second); err == nil {
		_ = c2.Close()
		t.Errorf("t=%v: new connection after wait succeeded", at())
	} else {
		t.Logf("t=%v: new connection after wait blocked: %v", at(), err)
	}

	// 4. Is the connection opened in the window still usable?
	if status, err := httpHead(conn); err == nil {
		t.Errorf("t=%v: established connection still usable after policy applied: %s", at(), status)
	} else {
		t.Logf("t=%v: established connection cut: %v", at(), err)
	}
}
