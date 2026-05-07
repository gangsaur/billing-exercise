package middleware

import (
	"log"
	"net/http"

	"gangsaur.com/billing-exercise/internal/static"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s - %s", r.Method, r.URL.Path, r.Context().Value(static.RequestIdKey))
		next.ServeHTTP(w, r)
	})
}
