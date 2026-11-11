package rmqenv

import (
	"fmt"
	"net/url"
)

// PostgresDSN returns postgres://user:password@localhost:5432/database?sslmode=disable.
// The credentials come from POSTGRES_USER / POSTGRES_PASSWORD / POSTGRES_DB in the same
// repository-root .env file (the database is reached through
// `kubectl port-forward svc/postgres 5432:5432`, see infra/k8s/postgres/README.md).
func PostgresDSN() (string, error) {
	vars, err := readEnvFile()
	if err != nil {
		return "", err
	}

	user, pass, db := vars["POSTGRES_USER"], vars["POSTGRES_PASSWORD"], vars["POSTGRES_DB"]
	if user == "" || pass == "" || db == "" {
		return "", fmt.Errorf("POSTGRES_USER / POSTGRES_PASSWORD / POSTGRES_DB are not set in .env")
	}

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(user, pass),
		Host:     "localhost:5432",
		Path:     "/" + db,
		RawQuery: "sslmode=disable",
	}
	return u.String(), nil
}
