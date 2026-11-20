package observability

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// AccessLog is middleware that writes one structured log line per request
// with method, path (without query string), status, duration_ms and request_id.
func AccessLog(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			completed := false
			defer func() {
				status := ww.Status()
				if status == 0 {
					if !completed {
						status = http.StatusInternalServerError
					} else {
						status = http.StatusOK
					}
				}

				durationMS := time.Since(start).Milliseconds()

				logger.InfoContext(r.Context(), "http request",
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Int("status", status),
					slog.Int64("duration_ms", durationMS),
					slog.String("request_id", RequestIDFromContext(r.Context())),
				)
			}()

			next.ServeHTTP(ww, r)
			completed = true
		})
	}
}
