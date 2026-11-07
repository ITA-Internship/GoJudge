// Package rmqenv builds the AMQP URI from RABBITMQ_USER / RABBITMQ_PASSWORD
// stored in the .env file in the repository root, so no secrets are kept in code.
package rmqenv

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// URI returns amqp://user:password@localhost:5672/.
func URI() (string, error) {
	vars, err := readEnvFile()
	if err != nil {
		return "", err
	}

	user, pass := vars["RABBITMQ_USER"], vars["RABBITMQ_PASSWORD"]
	if user == "" || pass == "" {
		return "", fmt.Errorf("RABBITMQ_USER / RABBITMQ_PASSWORD are not set in .env")
	}

	u := url.URL{Scheme: "amqp", User: url.UserPassword(user, pass), Host: "localhost:5672", Path: "/"}
	return u.String(), nil
}

// readEnvFile parses KEY=VALUE lines from <repo root>/.env.
// The root is derived from this file's location (spike/spk-03/RabbitMQ/rmqenv),
// so it works regardless of the directory the program is started from.
func readEnvFile() (map[string]string, error) {
	_, thisFile, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "..")

	f, err := os.Open(filepath.Join(root, ".env"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	vars := map[string]string{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if k, v, ok := strings.Cut(line, "="); ok {
			vars[strings.TrimSpace(k)] = strings.TrimSpace(v)
		}
	}
	return vars, sc.Err()
}
