package main

import (
	"bytes"
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ITA-Internship/GoJudge/internal/httpserver"
)

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

func env(m map[string]string) httpserver.LookupFunc {
	return func(k string) (string, bool) { v, ok := m[k]; return v, ok }
}

func TestRun_InvalidConfigReturnsError(t *testing.T) {
	err := run(context.Background(), env(map[string]string{httpserver.EnvShutdownTimeout: "0s"}), &syncBuffer{})
	if err == nil || !strings.Contains(err.Error(), "ShutdownTimeout") {
		t.Fatalf("run error = %v, want ShutdownTimeout validation error", err)
	}
}

func TestRun_StopsCleanlyOnCancel(t *testing.T) {
	out := &syncBuffer{}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- run(ctx, env(map[string]string{httpserver.EnvAddr: "127.0.0.1:0"}), out) }()

	deadline := time.Now().Add(2 * time.Second)
	for !strings.Contains(out.String(), "http server started") {
		if time.Now().After(deadline) {
			t.Fatalf("server did not start: %s", out.String())
		}
		time.Sleep(10 * time.Millisecond)
	}
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("run returned %v, want nil", err)
		}
	case <-time.After(httpserver.DefaultShutdownTimeout):
		t.Fatal("run did not return within the shutdown deadline")
	}
	if !strings.Contains(out.String(), "http server stopped") {
		t.Errorf("missing stop log line: %s", out.String())
	}
}
