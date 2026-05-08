package psql

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (p *Psql) GetLoan(ctx context.Context, id int) (Loan, error) {
	sql := "SELECT id, duration, principal_amount, outstanding_amount, interest, user_id, created_at, updated_at FROM loans WHERE id = $1"

	var loan Loan
	err := p.pool.QueryRow(ctx, sql, id).Scan(
		&loan.Id, &loan.Duration, &loan.PrincipalAmount, &loan.OutstandingAmount, &loan.InterestRate, &loan.UserId, &loan.CreatedAt, &loan.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return Loan{}, ErrNotFound
	}

	return loan, err
}
