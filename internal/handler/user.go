package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gangsaur.com/billing-exercise/internal/db/psql"
	"gangsaur.com/billing-exercise/internal/service"
)

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
