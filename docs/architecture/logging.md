# Logging and request correlation

Package: `internal/observability`. It uses only the standard library (`log/slog`, `net/http`, `crypto/rand`).

## What the package provides

| Symbol | Purpose |
|---|---|
| `NewLogger(w io.Writer, level slog.Level) *slog.Logger` | JSON logger, one object per line. |
| `RequestID(next http.Handler) http.Handler` | Middleware that gives every request an ID. |
| `RequestIDFromContext(ctx) string` | Returns the request ID, or `""` when none is set. |
| `AccessLog(logger) func(http.Handler) http.Handler` | Middleware that writes one log line per request. |
| `RequestIDHeader` | The header name, `X-Request-ID`. |

The middleware are plain `http.Handler` wrappers. They do not choose a router or an HTTP server.

## Log format

Every line is one JSON object with `time`, `level` and `msg`, plus structured fields.

Example access-log line:

```json
{"time":"2026-11-20T10:15:30.123456789Z","level":"INFO","msg":"http request","method":"GET","path":"/x","status":200,"duration_ms":3,"request_id":"3f2a9c0d5e7b41a8b6c1d2e3f4a5b6c7"}
```

## Allowed log fields

| Field | Type | Meaning |
|---|---|---|
| `time` | string | Timestamp added by `slog`. |
| `level` | string | `DEBUG`, `INFO`, `WARN` or `ERROR`. |
| `msg` | string | Short, fixed event description. |
| `method` | string | HTTP method. |
| `path` | string | URL path only, never the query string. |
| `status` | int | HTTP response status code. |
| `duration_ms` | int | Time spent handling the request, in milliseconds. |
| `request_id` | string | The request ID (see below). |

Handlers may add other structured fields (for example `component`), provided they follow the rules below. Add any new field to this table in the same pull request.

## Must never be logged

- Credentials, `Authorization` headers, cookies, tokens, passwords or secrets. Headers are not logged.
- Learner-submitted source code. Request bodies are not logged.
- Query strings, because they can carry tokens.
- Personal data beyond the request ID, method, path, status and duration.

Do not put these in `msg` or in custom fields either. Do not pass whole `*http.Request` values, header maps or bodies to the logger.

## Request ID rule

1. If the inbound `X-Request-ID` is **1 to 64 characters from `[A-Za-z0-9_-]`**, it is kept unchanged.
2. If it is missing, empty, longer than 64 characters, or contains any other character, it is **replaced** by a generated ID: 32 lowercase hexadecimal characters (16 random bytes from `crypto/rand`).
3. The final ID is set on the **response header** `X-Request-ID`, stored in the **request context**, and written as `request_id` in the access-log line.

The ID is for correlation only. It is not a secret and must not be used for authentication or authorization.

## Middleware order

`RequestID` must be the outer wrapper and `AccessLog` the inner one:

```go
handler := observability.RequestID(observability.AccessLog(logger)(mux))
```

If the order is reversed, `AccessLog` runs before the ID exists and logs `request_id` as `""`.

## Using the logger and request ID in a handler

The logger is created once at startup and passed to handlers explicitly. The request ID is read from the request context and added to any log line the handler writes, so those lines can be matched with the access-log line.

```go
func NewSubmitHandler(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := observability.RequestIDFromContext(r.Context())

		logger.InfoContext(r.Context(), "submission accepted",
			slog.String("request_id", id),
		)

		w.WriteHeader(http.StatusAccepted)
	})
}
```

Creating the logger at startup:

```go
logger := observability.NewLogger(os.Stdout, slog.LevelInfo)
```

Records below the configured level are dropped.

## Error behaviour

- The middleware never writes a response body, so no internal detail reaches clients.
- If `crypto/rand` fails, `RequestID` answers with a bare `500` and an empty body, and the next handler is not called. Because `RequestID` is the outer wrapper, no access-log line is written for that request.

## Tests

Unit tests live in `internal/observability/observability_test.go` and run with `go test -race ./...`.
