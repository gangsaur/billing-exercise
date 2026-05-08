package service

import (
	"context"

	"gangsaur.com/billing-exercise/internal/repository/db/psql"
)

type LoanService struct {
	store Store
}

func NewLoanService(l Store) *LoanService {
	return &LoanService{
		store: l,
	}
}

func (l *LoanService) GetLoan(ctx context.Context, id int) (psql.Loan, error) {
	loanData, err := l.store.GetLoan(ctx, id)
	if err != nil {
		return psql.Loan{}, err
	}

	return loanData, nil
}
