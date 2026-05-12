package service

import (
	"context"
	"time"

	"gangsaur.com/billing-exercise/internal/repository/db/psql"
	"github.com/jackc/pgx/v5"
)

type Store interface {
	GetLoan(ctx context.Context, id int) (psql.Loan, error)
	GetLoanByUserIdAndStatus(ctx context.Context, userId int, status int) ([]psql.Loan, error)
	GetLoanAndLoanPayments(ctx context.Context, id int) (psql.Loan, []psql.LoanPayment, error)
	GetLoanAndLoanPaymentsByStatusDueDate(ctx context.Context, id int, status int, date time.Time, dueDateBeforeDate bool) (psql.Loan, []psql.LoanPayment, error)

	GetLoanPaymentsByLoanIdsStatusDueDate(ctx context.Context, loandIds []int, status int, date time.Time, dueDateBeforeDate bool) ([]psql.LoanPayment, error)
	ReduceLoanOutstandingAmountTx(ctx context.Context, tx pgx.Tx, outstandingDeduction, loanId, previousOutstanding int) error
	ReduceLoanOutstandingAmountStatusPaidTx(ctx context.Context, tx pgx.Tx, outstandingDeduction, loanId, previousOutstanding int) error
	UpdateLoanPaymentStatusPaidTx(ctx context.Context, tx pgx.Tx, loanPaymentIds []int) error

	GetUser(ctx context.Context, id int) (psql.User, error)

	Begin(ctx context.Context) (pgx.Tx, error)
	Commit(ctx context.Context, tx pgx.Tx) error
	Rollback(ctx context.Context, tx pgx.Tx) error
}
