package psql

import (
	"context"
	"database/sql"
	"time"
)

const LoanPaymentStatusScheduled = 0
const LoanPaymentStatusPaid = 1

type LoanPayment struct {
	Id        int
	Period    int
	Amount    int
	DueDate   time.Time
	PaidAt    time.Time
	Status    int
	LoanId    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Psql) GetLoanPaymentsByLoanIdsStatusDueDate(ctx context.Context, loandIds []int, status int, date time.Time, dueDateBeforeDate bool) ([]LoanPayment, error) {
	var q string
	if dueDateBeforeDate {
		q = `
			SELECT id, period, amount, due_date, paid_at, status, loan_id, created_at, updated_at FROM loan_payments
			WHERE loan_id = ANY($1) AND status = $2 AND due_date < $3;`
	} else {
		q = `
			SELECT id, period, amount, due_date, paid_at, status, loan_id, created_at, updated_at FROM loan_payments
			WHERE loan_id = ANY($1) AND status = $2 AND due_date > $3;`
	}

	rows, err := p.pool.Query(ctx, q, loandIds, status, date)
	if err != nil {
		return []LoanPayment{}, err
	}
	defer rows.Close()

	loanPayments := make([]LoanPayment, 0)
	for rows.Next() {
		var loanPayment LoanPayment
		var paidAt sql.NullTime
		err := rows.Scan(&loanPayment.Id, &loanPayment.Period, &loanPayment.Amount, &loanPayment.DueDate, &paidAt, &loanPayment.Status, &loanPayment.LoanId, &loanPayment.CreatedAt, &loanPayment.UpdatedAt)
		if err != nil {
			return []LoanPayment{}, err
		}

		if paidAt.Valid {
			loanPayment.PaidAt = paidAt.Time
		}

		loanPayments = append(loanPayments, loanPayment)
	}

	return loanPayments, nil
}
