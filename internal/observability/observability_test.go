package observability

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

var hexIDPattern = regexp.MustCompile(`^[0-9a-f]{32}$`)

func TestNewLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, slog.LevelInfo)

	logger.Debug("ignored debug message")
	logger.Info("first event", slog.String("component", "api"))
	logger.Error("second event", slog.Int("code", 42))

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 log lines, got %d: %q", len(lines), buf.String())
	}

	for i, line := range lines {
		var entry map[string]any
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			t.Fatalf("line %d is not valid JSON (%q): %v", i, line, err)
		}
		for _, required := range []string{"time", "level", "msg"} {
			if _, ok := entry[required]; !ok {
				t.Errorf("line %d missing required field %q: %v", i, required, entry)
			}
		}
	}

	var first map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &first); err != nil {
		t.Fatalf("unmarshal first line: %v", err)
	}
	if first["level"] != "INFO" || first["msg"] != "first event" || first["component"] != "api" {
		t.Errorf("unexpected first log entry: %v", first)
	}
}

func TestRequestIDFromContext_Empty(t *testing.T) {
	if got := RequestIDFromContext(context.Background()); got != "" {
		t.Fatalf("expected empty string when middleware did not run, got %q", got)
	}
}

func TestRequestID_ValidPropagation(t *testing.T) {
	const inboundID = "Abc_123-xyz_64charsMaxCheck"
	var buf bytes.Buffer
	logger := NewLogger(&buf, slog.LevelInfo)

	var ctxID string
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxID = RequestIDFromContext(r.Context())
		w.WriteHeader(http.StatusAccepted)
	})

	handler := RequestID(AccessLog(logger)(inner))

	req := httptest.NewRequest(http.MethodPost, "/submit", nil)
	req.Header.Set(RequestIDHeader, inboundID)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get(RequestIDHeader); got != inboundID {
		t.Errorf("response header X-Request-ID = %q, want %q", got, inboundID)
	}
	if ctxID != inboundID {
		t.Errorf("RequestIDFromContext = %q, want %q", ctxID, inboundID)
	}

	var entry map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &entry); err != nil {
		t.Fatalf("unmarshal access log: %v", err)
	}
	if entry["request_id"] != inboundID {
		t.Errorf("access log request_id = %v, want %q", entry["request_id"], inboundID)
	}
}

func TestRequestID_InvalidOrMissingReplaced(t *testing.T) {
	cases := []struct {
		name  string
		value string
		set   bool
	}{
		{name: "missing", set: false},
		{name: "empty", value: "", set: true},
		{name: "too_long_65", value: strings.Repeat("a", 65), set: true},
		{name: "contains_space", value: "abc def", set: true},
		{name: "contains_slash", value: "abc/def", set: true},
		{name: "contains_symbol", value: "abc@def", set: true},
	}

	seen := make(map[string]struct{})

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var ctxID string
			handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctxID = RequestIDFromContext(r.Context())
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			if tc.set {
				req.Header.Set(RequestIDHeader, tc.value)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			respID := rec.Header().Get(RequestIDHeader)
			if !hexIDPattern.MatchString(respID) {
				t.Fatalf("generated ID %q does not match 32-char lowercase hex", respID)
			}
			if respID != ctxID {
				t.Errorf("header ID %q != context ID %q", respID, ctxID)
			}
			if _, exists := seen[respID]; exists {
				t.Errorf("duplicate generated ID %q", respID)
			}
			seen[respID] = struct{}{}
		})
	}
}

func TestRequestID_EntropyFailure(t *testing.T) {
	origRandRead := randRead
	t.Cleanup(func() { randRead = origRandRead })

	randRead = func([]byte) (int, error) {
		return 0, errors.New("entropy exhausted")
	}

	called := false
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if called {
		t.Fatal("next handler should not be called on entropy failure")
	}
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("expected empty body on 500, got %q", rec.Body.String())
	}
}

func TestAccessLog_StripsQueryString(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, slog.LevelInfo)

	handler := RequestID(AccessLog(logger)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))

	req := httptest.NewRequest(http.MethodGet, "/x?token=SENTINEL-1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	rawLog := buf.String()
	if strings.Contains(rawLog, "SENTINEL-1") {
		t.Fatalf("access log leaked query parameter: %s", rawLog)
	}

	var entry map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &entry); err != nil {
		t.Fatalf("unmarshal access log: %v", err)
	}
	if entry["path"] != "/x" {
		t.Errorf("path = %v, want %q", entry["path"], "/x")
	}
	if entry["method"] != http.MethodGet {
		t.Errorf("method = %v, want %q", entry["method"], http.MethodGet)
	}
	if entry["status"] != float64(http.StatusOK) {
		t.Errorf("status = %v, want %d", entry["status"], http.StatusOK)
	}
	if _, ok := entry["duration_ms"]; !ok {
		t.Errorf("missing duration_ms in log: %v", entry)
	}
	if id, ok := entry["request_id"].(string); !ok || !hexIDPattern.MatchString(id) {
		t.Errorf("invalid request_id in log: %v", entry["request_id"])
	}
}

func TestAccessLog_NeverLogsSecretsOrBody(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, slog.LevelInfo)

	handler := RequestID(AccessLog(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		w.WriteHeader(http.StatusNoContent)
	})))

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("payload=SENTINEL-3"))
	req.Header.Set("Authorization", "Bearer SENTINEL-2")
	req.AddCookie(&http.Cookie{
		Name:     "session",
		Value:    "COOKIE-SECRET",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	rawLog := buf.String()
	for _, secret := range []string{"SENTINEL-2", "SENTINEL-3", "COOKIE-SECRET"} {
		if strings.Contains(rawLog, secret) {
			t.Fatalf("access log leaked %q: %s", secret, rawLog)
		}
	}
}

func TestMiddleware_BodyPassthroughAndBare500(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, slog.LevelInfo)

	t.Run("200_passthrough", func(t *testing.T) {
		handler := RequestID(AccessLog(logger)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, err := w.Write([]byte("expected-payload"))
			if err != nil {
				t.Errorf("write failed: %v", err)
			}
		})))

		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want 200", rec.Code)
		}
		if rec.Body.String() != "expected-payload" {
			t.Errorf("body = %q, want %q", rec.Body.String(), "expected-payload")
		}
	})

	t.Run("500_adds_no_body", func(t *testing.T) {
		handler := RequestID(AccessLog(logger)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})))

		req := httptest.NewRequest(http.MethodGet, "/fail", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500", rec.Code)
		}
		if rec.Body.Len() != 0 {
			t.Errorf("expected empty body on 500, got %q", rec.Body.String())
		}
	})
}

func TestAccessLog_PanickingHandler(t *testing.T) {
	t.Run("panic_before_write_logs_500", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(&buf, slog.LevelInfo)

		handler := RequestID(AccessLog(logger)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			panic("boom")
		})))

		req := httptest.NewRequest(http.MethodGet, "/panic?token=SENTINEL-1", nil)
		req.Header.Set(RequestIDHeader, "panic-req-1")
		rec := httptest.NewRecorder()

		var recovered any
		func() {
			defer func() {
				recovered = recover()
			}()
			handler.ServeHTTP(rec, req)
		}()

		if recovered != "boom" {
			t.Fatalf("panic = %v, want %q", recovered, "boom")
		}
		if rec.Body.Len() != 0 {
			t.Errorf("expected empty body on panic, got %q", rec.Body.String())
		}
		if strings.Contains(buf.String(), "SENTINEL-1") {
			t.Fatalf("access log leaked query parameter: %s", buf.String())
		}

		var entry map[string]any
		if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &entry); err != nil {
			t.Fatalf("unmarshal access log: %v", err)
		}
		if entry["method"] != http.MethodGet {
			t.Errorf("method = %v, want %q", entry["method"], http.MethodGet)
		}
		if entry["path"] != "/panic" {
			t.Errorf("path = %v, want %q", entry["path"], "/panic")
		}
		if entry["status"] != float64(http.StatusInternalServerError) {
			t.Errorf("status = %v, want %d", entry["status"], http.StatusInternalServerError)
		}
		if entry["request_id"] != "panic-req-1" {
			t.Errorf("request_id = %v, want %q", entry["request_id"], "panic-req-1")
		}
		if _, ok := entry["duration_ms"]; !ok {
			t.Errorf("missing duration_ms in log: %v", entry)
		}
	})

	t.Run("panic_after_write_keeps_written_status", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(&buf, slog.LevelInfo)

		handler := RequestID(AccessLog(logger)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusAccepted)
			panic("late-boom")
		})))

		req := httptest.NewRequest(http.MethodPost, "/partial", nil)
		rec := httptest.NewRecorder()

		var recovered any
		func() {
			defer func() {
				recovered = recover()
			}()
			handler.ServeHTTP(rec, req)
		}()

		if recovered != "late-boom" {
			t.Fatalf("panic = %v, want %q", recovered, "late-boom")
		}

		var entry map[string]any
		if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &entry); err != nil {
			t.Fatalf("unmarshal access log: %v", err)
		}
		if entry["status"] != float64(http.StatusAccepted) {
			t.Errorf("status = %v, want %d", entry["status"], http.StatusAccepted)
		}
	})
}