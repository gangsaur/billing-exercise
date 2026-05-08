package service

import (
	"context"

	"gangsaur.com/billing-exercise/internal/db/psql"
)

type UserStore interface {
	GetUser(ctx context.Context, id int) (psql.User, error)
}

type UserService struct {
	store UserStore
}

func NewUserService(u UserStore) *UserService {
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
