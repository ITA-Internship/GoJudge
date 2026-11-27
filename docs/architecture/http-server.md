# HTTP server

Package: `internal/httpserver`. Entry point: `cmd/gojudge`. It uses only the standard library (`net/http`, `log/slog`) and `internal/observability`.

## Starting and stopping

````bash
go run ./cmd/gojudge
````

The server logs `http server started` with the listen address. It stops on SIGINT (Ctrl+C) or SIGTERM:

1. Readiness switches to `503` at once.
2. The listener is closed. In-flight requests are allowed to finish.
3. If they do not finish within `GOJUDGE_HTTP_SHUTDOWN_TIMEOUT`, the remaining connections are closed and the process exits with code 1 and the error.

Invalid configuration or a listen failure (for example, the port is in use) stops startup with exit code 1.

## Configuration

Configuration is read from environment variables. An unset variable uses its default. A set variable, even an empty one, must be valid, otherwise startup fails and every invalid value is listed. Durations use Go syntax (`500ms`, `5s`, `1m`) and must be greater than zero.

| Variable | Default | Meaning |
|---|---|---|
| `GOJUDGE_HTTP_ADDR` | `:8080` | Listen address, `host:port`. `:8080` listens on all interfaces. |
| `GOJUDGE_HTTP_READ_HEADER_TIMEOUT` | `5s` | Time allowed to read request headers. Must not exceed the read timeout. |
| `GOJUDGE_HTTP_READ_TIMEOUT` | `15s` | Time allowed to read the whole request. |
| `GOJUDGE_HTTP_WRITE_TIMEOUT` | `15s` | Time allowed to write the response. |
| `GOJUDGE_HTTP_IDLE_TIMEOUT` | `60s` | How long an idle keep-alive connection is kept. |
| `GOJUDGE_HTTP_SHUTDOWN_TIMEOUT` | `10s` | Graceful-shutdown deadline. |

No configuration value is a secret.

## Health endpoints

| Endpoint | Status | Body | Meaning |
|---|---|---|---|
| `GET /livez` | `200` | `{"status":"ok"}` | The process is running. Never checks dependencies. |
| `GET /readyz` | `200` | `{"status":"ready"}` | The server is accepting requests. |
| `GET /readyz` | `503` | `{"status":"not_ready"}` | Not serving yet, or shutting down. |

Readiness reflects the HTTP server lifecycle **only**. It does not check PostgreSQL, Auth0, RabbitMQ or the judge, and a `200` does not mean they are healthy. Dependency checks need their own agreed scope.

Responses carry `Content-Type: application/json` and `Cache-Control: no-store`. Other methods return `405`.

## Middleware

Every handler, including the health endpoints, runs inside the INF-07 chain in the documented order ([logging.md](logging.md#middleware-order)):

````go
observability.RequestID(observability.AccessLog(logger)(mux))
````

So every health request gets an `X-Request-ID` response header and one access-log line with the same `request_id` and the status. Query strings, headers, cookies and bodies are not logged.

## Adding endpoints

Register handlers on the server before it starts. Patterns use `net/http` ServeMux syntax:

````go
srv, err := httpserver.New(cfg, logger)
// ...
srv.Handle("GET /api/v1/me", meHandler)
````

Registered handlers get the same middleware. `/livez` and `/readyz` are reserved. A conflicting pattern panics at startup.

## Tests

`internal/httpserver/*_test.go` and `cmd/gojudge/main_test.go`, run with `go test -race ./...`.
