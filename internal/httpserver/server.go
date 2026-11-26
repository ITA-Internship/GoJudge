// Package httpserver provides the GoJudge HTTP server: configuration,
// health endpoints, the INF-07 middleware chain and graceful shutdown.
package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/ITA-Internship/GoJudge/internal/observability"
)

// Server wraps http.Server with readiness state and a handler registry.
type Server struct {
	cfg     Config
	mux     *http.ServeMux
	httpSrv *http.Server
	ready   atomic.Bool

	shutdownOnce sync.Once
	shutdownErr  error
}

// New validates cfg and builds a server with the health endpoints registered
// and the middleware chain RequestID(AccessLog(mux)) in place.
func New(cfg Config, logger *slog.Logger) (*Server, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if logger == nil {
		return nil, errors.New("httpserver: logger must not be nil")
	}

	s := &Server{cfg: cfg, mux: http.NewServeMux()}

	s.mux.Handle("GET "+LivenessPath, livenessHandler())
	s.mux.Handle("GET "+ReadinessPath, readinessHandler(&s.ready))

	// Documented order (docs/architecture/logging.md): RequestID outside AccessLog.
	handler := observability.RequestID(observability.AccessLog(logger)(s.mux))

	s.httpSrv = &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
	return s, nil
}

// Handle registers a handler for future endpoints (authentication,
// functional API). Call it before Serve. Patterns use the net/http
// ServeMux syntax, for example "GET /api/v1/me". Registered handlers run
// inside the same middleware chain as the health endpoints. Like
// ServeMux.Handle, it panics on a conflicting pattern.
func (s *Server) Handle(pattern string, h http.Handler) {
	s.mux.Handle(pattern, h)
}

// Handler returns the full handler chain, for tests and embedding.
func (s *Server) Handler() http.Handler { return s.httpSrv.Handler }

// Ready reports whether readiness currently returns 200.
func (s *Server) Ready() bool { return s.ready.Load() }

// Serve accepts connections on ln until Shutdown. It marks the server ready
// while serving. It returns nil after a graceful shutdown and the
// underlying error otherwise.
func (s *Server) Serve(ln net.Listener) error {
	s.ready.Store(true)
	defer s.ready.Store(false)

	err := s.httpSrv.Serve(ln)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return fmt.Errorf("http server: %w", err)
}

// Shutdown marks the server not ready, stops accepting connections and
// waits for in-flight requests, bounded by Config.ShutdownTimeout (or ctx,
// whichever ends first). If the deadline passes, remaining connections are
// closed and the deadline error is returned. Repeated calls return the
// first result.
func (s *Server) Shutdown(ctx context.Context) error {
	s.shutdownOnce.Do(func() {
		s.ready.Store(false)

		ctx, cancel := context.WithTimeout(ctx, s.cfg.ShutdownTimeout)
		defer cancel()

		if err := s.httpSrv.Shutdown(ctx); err != nil {
			closeErr := s.httpSrv.Close()
			s.shutdownErr = fmt.Errorf("http server shutdown: %w", errors.Join(err, closeErr))
		}
	})
	return s.shutdownErr
}
