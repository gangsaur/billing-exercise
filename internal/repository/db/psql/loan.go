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
		&loan.Id, &loan.Duration, &loan.PrincipalAmount, &loan.OutstandingAmount, &loan.Status, &loan.InterestRate, &loan.UserId, &loan.CreatedAt, &loan.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return Loan{}, ErrNotFound
	}

	return loan, err
}

func (p *Psql) GetLoanAndLoanPaymentsByStatusDueDate(ctx context.Context, id int, status int, date time.Time, dueDateBeforeDate bool) (Loan, []LoanPayment, error) {
	var q string

	if dueDateBeforeDate {
		q = `
			SELECT l.id, l.interest, l.outstanding_amount, lp.id, lp.amount, lp.loan_id FROM loans AS l JOIN loan_payments AS lp ON l.id = lp.loan_id
			WHERE l.id = $1 AND lp.status = $2 AND lp.due_date < $3;`
	} else {
		q = `
			SELECT l.id, l.interest, l.outstanding_amount, lp.id, lp.amount, lp.loan_id FROM loans AS l JOIN loan_payments AS lp ON l.id = lp.loan_id
			WHERE l.id = $1 AND lp.status = $2 AND lp.due_date > $3;`

	}

	rows, err := p.pool.Query(ctx, q, id, status, date)
	if err != nil {
		return Loan{}, []LoanPayment{}, err
	}
	defer rows.Close()

	// Scan 1st row for loan information, then scan the rest
	var loanPayments []LoanPayment
	var loan Loan
	if rows.Next() {
		var loanPayment LoanPayment
		err := rows.Scan(&loan.Id, &loan.InterestRate, &loan.OutstandingAmount, &loanPayment.Id, &loanPayment.Amount, &loanPayment.LoanId)
		if err != nil {
			return Loan{}, []LoanPayment{}, err
		}

		loanPayments = append(loanPayments, loanPayment)
	}

	for rows.Next() {
		var loanPayment LoanPayment
		err := rows.Scan(nil, nil, nil, &loanPayment.Id, &loanPayment.Amount, &loanPayment.LoanId)
		if err != nil {
			return Loan{}, []LoanPayment{}, err
		}

		loanPayments = append(loanPayments, loanPayment)
	}

	return loan, loanPayments, nil
}

func (p *Psql) GetLoanByUserIdAndStatus(ctx context.Context, userId int, status int) ([]Loan, error) {
	sql := "SELECT id, duration, principal_amount, outstanding_amount, status, interest, user_id, created_at, updated_at FROM loans WHERE user_id = $1 AND status = $2"

	rows, err := p.pool.Query(ctx, sql, userId, status)
	if err != nil {
		return []Loan{}, err
	}
	defer rows.Close()

	loans := make([]Loan, 0)
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

func (p *Psql) ReduceLoanOutstandingAmountTx(ctx context.Context, tx pgx.Tx, outstandingDeduction, loanId, previousOutstanding int) error {
	_, err := tx.Exec(ctx, "UPDATE loans SET outstanding_amount = outstanding_amount - $1, updated_at = NOW() WHERE id = $2 AND outstanding_amount = $3",
		outstandingDeduction, loanId, previousOutstanding)
	return err
}

func (p *Psql) ReduceLoanOutstandingAmountStatusPaidTx(ctx context.Context, tx pgx.Tx, outstandingDeduction, loanId, previousOutstanding int) error {
	_, err := tx.Exec(ctx, "UPDATE loans SET outstanding_amount = outstanding_amount - $1, status = 1, updated_at = NOW() WHERE id = $2 AND outstanding_amount = $3",
		outstandingDeduction, loanId, previousOutstanding)
	return err
}

func (p *Psql) UpdateLoanPaymentStatusPaidTx(ctx context.Context, tx pgx.Tx, loanPaymentIds []int) error {
	_, err := tx.Exec(ctx, "UPDATE loan_payments SET status = 1, paid_at = NOW(), updated_at = NOW() WHERE id = ANY($1) AND status=0", loanPaymentIds)
	return err
}
