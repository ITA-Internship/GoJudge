// Package observability provides structured logging and request correlation
// for HTTP handlers.
//
// The middleware are plain http.Handler wrappers, so they work with any router.
package observability

import (
	"io"
	"log/slog"
)

// NewLogger returns a logger that writes one JSON object per line to w.
// Every line has "time", "level" and "msg" plus any structured fields.
// Records below level are dropped.
func NewLogger(w io.Writer, level slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level}))
}