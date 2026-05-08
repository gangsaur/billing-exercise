package service

import (
	"context"

	"gangsaur.com/billing-exercise/internal/db/psql"
)

type Store interface {
	GetLoan(ctx context.Context, id int) (psql.Loan, error)
	GetUser(ctx context.Context, id int) (psql.User, error)
}
