package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gangsaur.com/billing-exercise/internal/db/psql"
	"gangsaur.com/billing-exercise/internal/service"
)

type UserResponse struct {
	Id         int  `json:"id"`
	Delinquent bool `json:"delinquent,omitempty"`
}

func toUserResponse(user psql.User) UserResponse {
	return UserResponse{
		Id: user.Id,
	}
}

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(u *service.UserService) *UserHandler {
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

		user, err := u.userService.GetUser(r.Context(), id)
		if err != nil {
			WriteErrorResponse(r.Context(), w, r, err)
			return
		}

		res, err := json.Marshal(toUserResponse(user))
		if err != nil {
			WriteErrorResponse(r.Context(), w, r, err)
			return
		}

		w.Write(res)
	}
}
