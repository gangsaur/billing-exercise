package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"gangsaur.com/billing-exercise/internal/repository/db/psql"
	"gangsaur.com/billing-exercise/internal/static"
)

func WriteErrorResponse(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("%s - %s", err.Error(), r.Context().Value(static.RequestIdKey))

	if errors.Is(err, psql.ErrNotFound) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write([]byte("{}"))
}
