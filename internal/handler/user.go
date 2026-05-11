package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"gangsaur.com/billing-exercise/internal/repository/db/psql"
)

// Response

type UserResponseDelinquent struct {
	Id         int  `json:"id"`
	Delinquent bool `json:"delinquent"`
}

func toUserResponseDelinquent(user psql.User, delinquentStatus bool) UserResponseDelinquent {
	return UserResponseDelinquent{
		Id:         user.Id,
		Delinquent: delinquentStatus,
	}
}

// Interface

type UserService interface {
	GetUser(ctx context.Context, id int) (psql.User, bool, error)
}

// Handler

type UserHandler struct {
	userService UserService
}

func NewUserHandler(u UserService) *UserHandler {
	return &UserHandler{
		userService: u,
	}
}

func (u *UserHandler) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			WriteErrorResponse(r.Context(), w, r, err)
			return
		}

		user, delinquentStatus, err := u.userService.GetUser(r.Context(), id)
		if err != nil {
			WriteErrorResponse(r.Context(), w, r, err)
			return
		}

		res, err := json.Marshal(toUserResponseDelinquent(user, delinquentStatus))
		if err != nil {
			WriteErrorResponse(r.Context(), w, r, err)
			return
		}

		w.Write(res)
	}
}
