package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ITA-Internship/GoJudge/internal/observability"
)

func discardLogger() *slog.Logger { return slog.New(slog.NewJSONHandler(io.Discard, nil)) }

// syncBuffer lets the server goroutines and the test share a log buffer.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

func testConfig() Config {
	cfg := DefaultConfig()
	cfg.Addr = "127.0.0.1:0"
	cfg.ShutdownTimeout = 2 * time.Second
	return cfg
}

// startServer runs s on a random loopback port and returns its base URL and
// a channel with Serve's result.
func startServer(t *testing.T, s *Server) (string, <-chan error) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	done := make(chan error, 1)
	go func() { done <- s.Serve(ln) }()
	t.Cleanup(func() { _ = s.Shutdown(context.Background()) })
	return "http://" + ln.Addr().String(), done
}

func get(t *testing.T, url string, header http.Header) (*http.Response, string) {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	for k, v := range header {
		req.Header[k] = v
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, _ := io.ReadAll(resp.Body)
	return resp, string(body)
}

func TestHealth_WhileRunning(t *testing.T) {
	s, err := New(testConfig(), discardLogger())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	base, _ := startServer(t, s)

	resp, body := get(t, base+LivenessPath, nil)
	if resp.StatusCode != http.StatusOK || !strings.Contains(body, `"ok"`) {
		t.Fatalf("liveness = %d %q, want 200 ok", resp.StatusCode, body)
	}

	resp, body = get(t, base+ReadinessPath, nil)
	if resp.StatusCode != http.StatusOK || !strings.Contains(body, `"ready"`) {
		t.Fatalf("readiness = %d %q, want 200 ready", resp.StatusCode, body)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q", ct)
	}
}

func TestReadiness_NotReadyBeforeServe(t *testing.T) {
	s, err := New(testConfig(), discardLogger())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, ReadinessPath, nil))
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("readiness before Serve = %d, want 503", rec.Code)
	}
}

func TestHealth_MiddlewareOrderAndAccessLog(t *testing.T) {
	logs := &syncBuffer{}
	s, err := New(testConfig(), observability.NewLogger(logs, slog.LevelInfo))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	base, _ := startServer(t, s)

	t.Run("generated id", func(t *testing.T) {
		resp, _ := get(t, base+ReadinessPath+"?token=QUERY-SENTINEL", nil)
		id := resp.Header.Get(observability.RequestIDHeader)
		if id == "" {
			t.Fatal("missing X-Request-ID response header")
		}
		entry := findAccessLog(t, logs.String(), id)
		if entry["status"] != float64(http.StatusOK) {
			t.Errorf("status = %v, want 200", entry["status"])
		}
		if entry["path"] != ReadinessPath {
			t.Errorf("path = %v, want %s", entry["path"], ReadinessPath)
		}
		if strings.Contains(logs.String(), "QUERY-SENTINEL") {
			t.Fatalf("query string leaked into access log: %s", logs.String())
		}
	})

	t.Run("inbound id kept", func(t *testing.T) {
		const inbound = "probe-123"
		resp, _ := get(t, base+LivenessPath, http.Header{observability.RequestIDHeader: {inbound}})
		if got := resp.Header.Get(observability.RequestIDHeader); got != inbound {
			t.Fatalf("X-Request-ID = %q, want %q", got, inbound)
		}
		// If AccessLog wrapped RequestID instead, request_id would be "".
		findAccessLog(t, logs.String(), inbound)
	})
}

// findAccessLog returns the access-log entry whose request_id is id.
func findAccessLog(t *testing.T, raw, id string) map[string]any {
	t.Helper()
	for _, line := range strings.Split(strings.TrimSpace(raw), "\n") {
		var e map[string]any
		if json.Unmarshal([]byte(line), &e) != nil {
			continue
		}
		if e["msg"] == "http request" && e["request_id"] == id {
			return e
		}
	}
	t.Fatalf("no access-log record with request_id %q in:\n%s", id, raw)
	return nil
}

func TestRegistrationPoint(t *testing.T) {
	logs := &syncBuffer{}
	s, err := New(testConfig(), observability.NewLogger(logs, slog.LevelInfo))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	s.Handle("GET /test/echo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, observability.RequestIDFromContext(r.Context()))
	}))
	base, _ := startServer(t, s)

	resp, body := get(t, base+"/test/echo", nil)
	if resp.StatusCode != http.StatusOK || body == "" || body != resp.Header.Get(observability.RequestIDHeader) {
		t.Fatalf("registered handler not behind the middleware chain: %d %q", resp.StatusCode, body)
	}
	findAccessLog(t, logs.String(), body)
}

func TestShutdown_GracefulWithinDeadline(t *testing.T) {
	cfg := testConfig()
	s, err := New(cfg, discardLogger())
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	started := make(chan struct{})
	release := make(chan struct{})
	s.Handle("GET /test/slow", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		close(started)
		<-release
		w.WriteHeader(http.StatusOK)
	}))
	base, serveDone := startServer(t, s)

	if resp, _ := get(t, base+ReadinessPath, nil); resp.StatusCode != http.StatusOK {
		t.Fatalf("readiness before shutdown = %d, want 200", resp.StatusCode)
	}

	inflight := make(chan int, 1)
	go func() {
		resp, err := http.Get(base + "/test/slow")
		if err != nil {
			inflight <- -1
			return
		}
		_ = resp.Body.Close()
		inflight <- resp.StatusCode
	}()
	<-started

	begin := time.Now()
	shutdownDone := make(chan error, 1)
	go func() { shutdownDone <- s.Shutdown(context.Background()) }()

	// During shutdown: not ready, both by flag and through the handler chain.
	waitFor(t, func() bool { return !s.Ready() })
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, ReadinessPath, nil))
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("readiness during shutdown = %d, want 503", rec.Code)
	}

	close(release)

	if err := <-shutdownDone; err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
	if elapsed := time.Since(begin); elapsed > cfg.ShutdownTimeout {
		t.Fatalf("shutdown took %s, deadline %s", elapsed, cfg.ShutdownTimeout)
	}
	if code := <-inflight; code != http.StatusOK {
		t.Fatalf("in-flight request = %d, want 200 (drained)", code)
	}
	if err := <-serveDone; err != nil {
		t.Fatalf("Serve returned %v after graceful shutdown, want nil", err)
	}
	if s.Ready() {
		t.Fatal("server still reports ready after shutdown")
	}
	if _, err := http.Get(base + LivenessPath); err == nil {
		t.Fatal("server still accepts connections after shutdown")
	}
}

func TestShutdown_DeadlineExceededIsSurfaced(t *testing.T) {
	cfg := testConfig()
	cfg.ShutdownTimeout = 100 * time.Millisecond
	s, err := New(cfg, discardLogger())
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	started := make(chan struct{})
	release := make(chan struct{})
	defer close(release)
	s.Handle("GET /test/stuck", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		close(started)
		<-release
	}))
	base, _ := startServer(t, s)

	go func() {
		if resp, err := http.Get(base + "/test/stuck"); err == nil {
			_ = resp.Body.Close()
		}
	}()
	<-started

	begin := time.Now()
	err = s.Shutdown(context.Background())
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Shutdown error = %v, want context.DeadlineExceeded", err)
	}
	if elapsed := time.Since(begin); elapsed > cfg.ShutdownTimeout+500*time.Millisecond {
		t.Fatalf("Shutdown took %s, deadline %s", elapsed, cfg.ShutdownTimeout)
	}
	if s.Ready() {
		t.Fatal("server reports ready after failed shutdown")
	}
}

func waitFor(t *testing.T, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatal("condition not met within 2s")
		}
		time.Sleep(5 * time.Millisecond)
	}
}
