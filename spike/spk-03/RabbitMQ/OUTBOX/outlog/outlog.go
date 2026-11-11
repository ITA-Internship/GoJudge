// Package outlog lets every SC-E program write its own log file.
// Open returns a writer that prints to the console and, if path is not empty,
// also to that file (created or truncated, parent folders are created).
// File writes are not buffered, so the log is complete even if the program is killed.
package outlog

import (
	"io"
	"os"
	"path/filepath"
)

func Open(path string) (io.Writer, func()) {
	if path == "" {
		return os.Stdout, func() {}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		panic(err)
	}
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	return io.MultiWriter(os.Stdout, f), func() { f.Close() }
}
