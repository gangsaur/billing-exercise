package middleware_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"gangsaur.com/billing-exercise/internal/middleware"
	"gangsaur.com/billing-exercise/internal/static"

	"github.com/stretchr/testify/assert"
)

func TestRequestMiddleware(t *testing.T) {
	tests := []struct {
		name            string
		next            http.Handler
		headerRequestId string
		want            []byte
	}{
		{
			name: "success with request ID set",
			next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(r.Context().Value(static.RequestIdKey).(string)))
			}),
			headerRequestId: "example-request-id",
			want:            []byte("example-request-id"),
		},
		{
			name: "success without request ID set",
			next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(r.Context().Value(static.RequestIdKey).(string)))
			}),
			headerRequestId: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/test", nil)
			r.Header.Add("X-Request-ID", tt.headerRequestId)
			w := httptest.NewRecorder()

			wrappedHandler := middleware.RequestMiddleware(tt.next)
			wrappedHandler.ServeHTTP(w, r)

			res := w.Result()
			body, _ := io.ReadAll(res.Body)
			res.Body.Close()

			if tt.headerRequestId != "" {
				assert.Equal(t, tt.want, bytes.TrimSpace(body))
			} else {
				match, _ := regexp.Match(`[a-z0-9]{8}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{12}`, body)
				assert.True(t, match, "Request ID not properly generated and set to context")
			}
		})
	}
}
