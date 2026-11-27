package httpserver

import (
	"net/http"
	"sync/atomic"
)

// Health endpoint paths. They sit outside /api/v1 because they are
// operational, not part of the public API.
const (
	LivenessPath  = "/livez"
	ReadinessPath = "/readyz"
)

var (
	bodyAlive    = []byte(`{"status":"ok"}` + "\n")
	bodyReady    = []byte(`{"status":"ready"}` + "\n")
	bodyNotReady = []byte(`{"status":"not_ready"}` + "\n")
)

// livenessHandler reports that the process is running. It never checks
// dependencies.
func livenessHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeHealth(w, http.StatusOK, bodyAlive)
	})
}

// readinessHandler reports the HTTP server lifecycle only: 200 while the
// server is accepting requests, 503 before Serve and once shutdown begins.
// It says nothing about the database, authentication, messaging or judge.
func readinessHandler(ready *atomic.Bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if ready.Load() {
			writeHealth(w, http.StatusOK, bodyReady)
			return
		}
		writeHealth(w, http.StatusServiceUnavailable, bodyNotReady)
	})
}

func writeHealth(w http.ResponseWriter, status int, body []byte) {
	h := w.Header()
	h.Set("Content-Type", "application/json")
	h.Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}
