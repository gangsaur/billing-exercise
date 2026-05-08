package handler

import (
	"context"
	"log"
	"net/http"

	"gangsaur.com/billing-exercise/internal/static"
)

func WriteGenericError(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("%s - %s", err.Error(), r.Context().Value(static.RequestIdKey))

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("{}"))
}
