package service

import (
	"context"
	"time"

	"gangsaur.com/billing-exercise/internal/repository/db/psql"
)

type Store interface {
	GetLoan(ctx context.Context, id int) (psql.Loan, error)
	GetLoanByUserIdAndStatus(ctx context.Context, userId int, status int) ([]psql.Loan, error)
	GetLoanPaymentsByLoanIdsStatusDueDate(ctx context.Context, loandIds []int, status int, date time.Time, dueDateBeforeDate bool) ([]psql.LoanPayment, error)
	GetUser(ctx context.Context, id int) (psql.User, error)
}
