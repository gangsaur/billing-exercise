package middleware_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"gangsaur.com/billing-exercise/internal/middleware"
	"gangsaur.com/billing-exercise/internal/static"

	"github.com/stretchr/testify/assert"
)

func TestLogMiddleware(t *testing.T) {
	tests := []struct {
		name    string
		next    http.Handler
		request *http.Request
		want    string
	}{
		{
			name:    "success",
			next:    http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			request: httptest.NewRequest("GET", "/test", nil).WithContext(context.WithValue(context.Background(), static.RequestIdKey, "example-request-id")),
			want:    "GET /test - example-request-id",
		},
	}

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stdout)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrappedHandler := middleware.LogMiddleware(tt.next)
			wrappedHandler.ServeHTTP(nil, tt.request)

			got := buf.String()
			check := strings.Contains(got, tt.want)
			assert.True(t, check, fmt.Sprintf("Incorrect logging format, want to contains %s, actual: %s", tt.want, got))
		})
	}
}
