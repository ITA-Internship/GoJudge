// Command gojudge runs the GoJudge HTTP server.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ITA-Internship/GoJudge/internal/httpserver"
	"github.com/ITA-Internship/GoJudge/internal/observability"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, os.LookupEnv, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "gojudge:", err)
		stop()
		os.Exit(1)
	}
}

// run starts the server and blocks until ctx is cancelled (SIGINT/SIGTERM)
// or the server fails. It returns any startup, serve or shutdown error.
func run(ctx context.Context, lookup httpserver.LookupFunc, out io.Writer) error {
	logger := observability.NewLogger(out, slog.LevelInfo)

	cfg, err := httpserver.LoadConfig(lookup)
	if err != nil {
		return err
	}
	srv, err := httpserver.New(cfg, logger)
	if err != nil {
		return err
	}

	// Future endpoints are registered here, before Serve, e.g.:
	// srv.Handle("GET /api/v1/me", meHandler)

	var lc net.ListenConfig
	ln, err := lc.Listen(ctx, "tcp", cfg.Addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", cfg.Addr, err)
	}

	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.Serve(ln) }()
	logger.Info("http server started", slog.String("addr", ln.Addr().String()))

	select {
	case err := <-serveErr:
		return err // Serve stopped on its own: a real failure.
	case <-ctx.Done():
	}

	logger.Info("http server shutting down")
	shutdownErr := srv.Shutdown(context.Background())
	err = errors.Join(shutdownErr, <-serveErr)
	if err == nil {
		logger.Info("http server stopped")
	}
	return err
}
