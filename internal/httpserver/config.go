package httpserver

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"
)

// Environment variables read by LoadConfig. An unset variable keeps its
// default. A set variable (even an empty one) must hold a valid value.
const (
	EnvAddr              = "GOJUDGE_HTTP_ADDR"
	EnvReadHeaderTimeout = "GOJUDGE_HTTP_READ_HEADER_TIMEOUT"
	EnvReadTimeout       = "GOJUDGE_HTTP_READ_TIMEOUT"
	EnvWriteTimeout      = "GOJUDGE_HTTP_WRITE_TIMEOUT"
	EnvIdleTimeout       = "GOJUDGE_HTTP_IDLE_TIMEOUT"
	EnvShutdownTimeout   = "GOJUDGE_HTTP_SHUTDOWN_TIMEOUT"
)

// Default values. Documented in docs/architecture/http-server.md.
const (
	DefaultAddr              = ":8080"
	DefaultReadHeaderTimeout = 5 * time.Second
	DefaultReadTimeout       = 15 * time.Second
	DefaultWriteTimeout      = 15 * time.Second
	DefaultIdleTimeout       = 60 * time.Second
	DefaultShutdownTimeout   = 10 * time.Second
)

// Config holds the HTTP server settings.
type Config struct {
	Addr              string        // listen address, host:port
	ReadHeaderTimeout time.Duration // http.Server.ReadHeaderTimeout
	ReadTimeout       time.Duration // http.Server.ReadTimeout
	WriteTimeout      time.Duration // http.Server.WriteTimeout
	IdleTimeout       time.Duration // http.Server.IdleTimeout
	ShutdownTimeout   time.Duration // graceful-shutdown deadline
}

// DefaultConfig returns the documented defaults.
func DefaultConfig() Config {
	return Config{
		Addr:              DefaultAddr,
		ReadHeaderTimeout: DefaultReadHeaderTimeout,
		ReadTimeout:       DefaultReadTimeout,
		WriteTimeout:      DefaultWriteTimeout,
		IdleTimeout:       DefaultIdleTimeout,
		ShutdownTimeout:   DefaultShutdownTimeout,
	}
}

// Validate reports every invalid field in one error.
func (c Config) Validate() error {
	var errs []error

	if err := validateAddr(c.Addr); err != nil {
		errs = append(errs, err)
	}

	durations := []struct {
		name string
		v    time.Duration
	}{
		{"ReadHeaderTimeout", c.ReadHeaderTimeout},
		{"ReadTimeout", c.ReadTimeout},
		{"WriteTimeout", c.WriteTimeout},
		{"IdleTimeout", c.IdleTimeout},
		{"ShutdownTimeout", c.ShutdownTimeout},
	}
	for _, d := range durations {
		if d.v <= 0 {
			errs = append(errs, fmt.Errorf("%s must be greater than zero, got %s", d.name, d.v))
		}
	}

	if c.ReadHeaderTimeout > 0 && c.ReadTimeout > 0 && c.ReadHeaderTimeout > c.ReadTimeout {
		errs = append(errs, fmt.Errorf("ReadHeaderTimeout (%s) must not exceed ReadTimeout (%s)",
			c.ReadHeaderTimeout, c.ReadTimeout))
	}

	if len(errs) > 0 {
		return fmt.Errorf("invalid http server config: %w", errors.Join(errs...))
	}
	return nil
}

func validateAddr(addr string) error {
	if addr == "" {
		return errors.New("addr must not be empty")
	}
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("addr %q is not host:port: %w", addr, err)
	}
	p, err := strconv.Atoi(port)
	if err != nil || p < 0 || p > 65535 {
		return fmt.Errorf("addr %q has an invalid port", addr)
	}
	return nil
}

// LookupFunc matches os.LookupEnv. Tests pass a map-backed function.
type LookupFunc func(key string) (string, bool)

// LoadConfig starts from DefaultConfig, applies any variables that are set,
// and validates the result.
func LoadConfig(lookup LookupFunc) (Config, error) {
	cfg := DefaultConfig()
	var errs []error

	if v, ok := lookup(EnvAddr); ok {
		cfg.Addr = v
	}

	durations := []struct {
		env string
		dst *time.Duration
	}{
		{EnvReadHeaderTimeout, &cfg.ReadHeaderTimeout},
		{EnvReadTimeout, &cfg.ReadTimeout},
		{EnvWriteTimeout, &cfg.WriteTimeout},
		{EnvIdleTimeout, &cfg.IdleTimeout},
		{EnvShutdownTimeout, &cfg.ShutdownTimeout},
	}
	for _, d := range durations {
		v, ok := lookup(d.env)
		if !ok {
			continue
		}
		parsed, err := time.ParseDuration(v)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", d.env, err))
			continue
		}
		*d.dst = parsed
	}

	if len(errs) > 0 {
		return Config{}, fmt.Errorf("invalid http server config: %w", errors.Join(errs...))
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
