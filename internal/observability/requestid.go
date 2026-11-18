package observability

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

// RequestIDHeader is the header used to receive and return the request ID.
const RequestIDHeader = "X-Request-ID"

const (
	maxRequestIDLen = 64
	generatedIDLen  = 16 // random bytes; 32 hex characters
)

type ctxKey struct{}

// randRead is replaced in tests to simulate an entropy failure.
var randRead = rand.Read

// RequestID is middleware that gives every request an ID.
//
// A valid inbound X-Request-ID (1-64 characters from [A-Za-z0-9_-]) is kept
// unchanged. A missing, empty, over-long or invalid one is replaced by a new
// 32-character lowercase hex ID from crypto/rand. The ID is set on the
// response header and stored in the request context.
//
// If random bytes cannot be read, the request fails with a bare 500 and no
// body, so no internal detail reaches the client.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(RequestIDHeader)
		if !validRequestID(id) {
			generated, err := newRequestID()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			id = generated
		}

		w.Header().Set(RequestIDHeader, id)
		ctx := context.WithValue(r.Context(), ctxKey{}, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestIDFromContext returns the request ID stored by RequestID, or ""
// when the context carries none.
func RequestIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(ctxKey{}).(string)
	return id
}

func validRequestID(id string) bool {
	if len(id) == 0 || len(id) > maxRequestIDLen {
		return false
	}
	for i := 0; i < len(id); i++ {
		c := id[i]
		switch {
		case c >= 'a' && c <= 'z':
		case c >= 'A' && c <= 'Z':
		case c >= '0' && c <= '9':
		case c == '_' || c == '-':
		default:
			return false
		}
	}
	return true
}

func newRequestID() (string, error) {
	b := make([]byte, generatedIDLen)
	if _, err := randRead(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
