package middleware

import (
	"net/http"
	"slices"
)

type MiddlewareChain []func(http.Handler) http.Handler

func (c MiddlewareChain) ThenFunc(h http.HandlerFunc) http.Handler {
	return c.Then(h)
}

func (c MiddlewareChain) Then(h http.Handler) http.Handler {
	for _, mw := range slices.Backward(c) {
		h = mw(h)
	}
	return h
}
