package middleware

import (
	"context"
	"net/http"

	"gangsaur.com/billing-exercise/internal/static"

	"github.com/google/uuid"
)

func RequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = uuid.NewString()
			r.Header.Set("X-Request-ID", id)
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), static.RequestIdKey, id)))
	})
}
