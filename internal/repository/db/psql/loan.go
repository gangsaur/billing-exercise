package psql

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

const LoanStatusOpen = 0
const LoanStatusClosed = 1

type Loan struct {
	Id                int
	Duration          int
	PrincipalAmount   int
	OutstandingAmount int
	Status            int
	InterestRate      float32
	UserId            int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (p *Psql) GetLoan(ctx context.Context, id int) (Loan, error) {
	sql := "SELECT id, duration, principal_amount, outstanding_amount, status, interest, user_id, created_at, updated_at FROM loans WHERE id = $1"

	var loan Loan
	err := p.pool.QueryRow(ctx, sql, id).Scan(
		&loan.Id, &loan.Duration, &loan.PrincipalAmount, &loan.OutstandingAmount, &loan.InterestRate, &loan.Status, &loan.UserId, &loan.CreatedAt, &loan.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return Loan{}, ErrNotFound
	}

	return loan, err
}

func (p *Psql) GetLoanByUserIdAndStatus(ctx context.Context, userId int, status int) ([]Loan, error) {
	sql := "SELECT id, duration, principal_amount, outstanding_amount, status, interest, user_id, created_at, updated_at FROM loans WHERE user_id = $1 AND status = $2"

	rows, err := p.pool.Query(ctx, sql, userId, status)
	if err != nil {
		return []Loan{}, err
	}
	defer rows.Close()

	loans := make([]Loan, 0, 0)
	for rows.Next() {
		var loan Loan
		err := rows.Scan(&loan.Id, &loan.Duration, &loan.PrincipalAmount, &loan.OutstandingAmount, &loan.InterestRate, &loan.Status, &loan.UserId, &loan.CreatedAt, &loan.UpdatedAt)
		if err != nil {
			return []Loan{}, err
		}

		loans = append(loans, loan)
	}

	return loans, nil
}
