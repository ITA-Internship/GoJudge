package httpserver

import (
	"strings"
	"testing"
	"time"
)

func mapLookup(m map[string]string) LookupFunc {
	return func(k string) (string, bool) {
		v, ok := m[k]
		return v, ok
	}
}

func TestLoadConfig_DefaultsWhenUnset(t *testing.T) {
	cfg, err := LoadConfig(mapLookup(nil))
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	want := Config{
		Addr:              ":8080",
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ShutdownTimeout:   10 * time.Second,
	}
	if cfg != want {
		t.Fatalf("cfg = %+v, want %+v", cfg, want)
	}
}

func TestLoadConfig_Overrides(t *testing.T) {
	cfg, err := LoadConfig(mapLookup(map[string]string{
		EnvAddr:              "127.0.0.1:9090",
		EnvReadHeaderTimeout: "2s",
		EnvReadTimeout:       "3s",
		EnvWriteTimeout:      "4s",
		EnvIdleTimeout:       "30s",
		EnvShutdownTimeout:   "1500ms",
	}))
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	want := Config{"127.0.0.1:9090", 2 * time.Second, 3 * time.Second, 4 * time.Second, 30 * time.Second, 1500 * time.Millisecond}
	if cfg != want {
		t.Fatalf("cfg = %+v, want %+v", cfg, want)
	}
}

func TestLoadConfig_Invalid(t *testing.T) {
	cases := map[string]struct {
		env     map[string]string
		wantErr string
	}{
		"bad duration":      {map[string]string{EnvReadTimeout: "fast"}, EnvReadTimeout},
		"empty duration":    {map[string]string{EnvIdleTimeout: ""}, EnvIdleTimeout},
		"zero timeout":      {map[string]string{EnvWriteTimeout: "0s"}, "WriteTimeout must be greater than zero"},
		"negative shutdown": {map[string]string{EnvShutdownTimeout: "-1s"}, "ShutdownTimeout must be greater than zero"},
		"empty addr":        {map[string]string{EnvAddr: ""}, "addr must not be empty"},
		"addr without port": {map[string]string{EnvAddr: "localhost"}, "not host:port"},
		"port out of range": {map[string]string{EnvAddr: ":70000"}, "invalid port"},
		"header > read":     {map[string]string{EnvReadHeaderTimeout: "20s"}, "must not exceed ReadTimeout"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := LoadConfig(mapLookup(tc.env))
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("error %q does not mention %q", err, tc.wantErr)
			}
		})
	}
}

func TestValidate_ReportsAllProblems(t *testing.T) {
	err := Config{}.Validate()
	if err == nil {
		t.Fatal("zero Config must be invalid")
	}
	for _, field := range []string{"addr", "ReadHeaderTimeout", "ReadTimeout", "WriteTimeout", "IdleTimeout", "ShutdownTimeout"} {
		if !strings.Contains(err.Error(), field) {
			t.Errorf("error does not mention %s: %v", field, err)
		}
	}
}

func TestNew_AppliesEffectiveTimeouts(t *testing.T) {
	cfg := DefaultConfig()
	s, err := New(cfg, discardLogger())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	got := s.httpSrv
	checks := []struct {
		name      string
		got, want time.Duration
	}{
		{"ReadHeaderTimeout", got.ReadHeaderTimeout, cfg.ReadHeaderTimeout},
		{"ReadTimeout", got.ReadTimeout, cfg.ReadTimeout},
		{"WriteTimeout", got.WriteTimeout, cfg.WriteTimeout},
		{"IdleTimeout", got.IdleTimeout, cfg.IdleTimeout},
		{"ShutdownTimeout", s.cfg.ShutdownTimeout, cfg.ShutdownTimeout},
	}
	for _, c := range checks {
		if c.got <= 0 {
			t.Errorf("%s = %s, want > 0", c.name, c.got)
		}
		if c.got != c.want {
			t.Errorf("%s = %s, want %s", c.name, c.got, c.want)
		}
	}
}

func TestNew_RejectsInvalidConfig(t *testing.T) {
	cfg := DefaultConfig()
	cfg.ShutdownTimeout = 0
	if _, err := New(cfg, discardLogger()); err == nil {
		t.Fatal("New accepted an invalid config")
	}
}
