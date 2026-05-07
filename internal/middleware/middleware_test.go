package middleware_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gangsaur.com/billing-exercise/internal/middleware"

	"github.com/stretchr/testify/assert"
)

var (
	middlewareOne   = func(next http.Handler) http.Handler { return testMiddlewareGenerator(next, "1") }
	middlewareTwo   = func(next http.Handler) http.Handler { return testMiddlewareGenerator(next, "2") }
	middlewareThree = func(next http.Handler) http.Handler { return testMiddlewareGenerator(next, "3") }
)

func testMiddlewareGenerator(next http.Handler, value string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { r.Header.Add("str", value); next.ServeHTTP(w, r) })
}

func TestMiddlewareChain_ThenFunc(t *testing.T) {
	tests := []struct {
		name          string
		startingChain []func(http.Handler) http.Handler
		h             http.HandlerFunc
		want          []byte
	}{
		{
			name:          "success",
			startingChain: []func(http.Handler) http.Handler{middlewareOne, middlewareThree, middlewareTwo},
			h: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(strings.Join(r.Header.Values("str"), "")))
			}),
			want: []byte("132"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			c := middleware.MiddlewareChain(tt.startingChain)
			completeChain := c.ThenFunc(tt.h)

			completeChain.ServeHTTP(w, r)

			res := w.Result()
			body, _ := io.ReadAll(res.Body)
			res.Body.Close()

			assert.Equal(t, tt.want, body)
		})
	}
}

func TestMiddlewareChain_Then(t *testing.T) {
	tests := []struct {
		name          string
		startingChain []func(http.Handler) http.Handler
		h             http.Handler
		want          []byte
	}{
		{
			name:          "success",
			startingChain: []func(http.Handler) http.Handler{middlewareOne, middlewareThree, middlewareTwo},
			h: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(strings.Join(r.Header.Values("str"), "")))
			}),
			want: []byte("132"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			c := middleware.MiddlewareChain(tt.startingChain)
			completeChain := c.Then(tt.h)

			completeChain.ServeHTTP(w, r)

			res := w.Result()
			body, _ := io.ReadAll(res.Body)
			res.Body.Close()

			assert.Equal(t, tt.want, body)
		})
	}
}
