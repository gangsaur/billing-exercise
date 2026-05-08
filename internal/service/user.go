package service

import (
	"context"

	"gangsaur.com/billing-exercise/internal/db/psql"
)

type UserService struct {
	store Store
}

func NewUserService(u Store) *UserService {
	return &UserService{
		store: u,
	}
}

func (u *UserService) GetUser(ctx context.Context, id int) (psql.User, error) {
	userData, err := u.store.GetUser(ctx, id)
	if err != nil {
		return psql.User{}, err
	}

	return userData, nil
}
