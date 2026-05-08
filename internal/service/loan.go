package service

import (
	"context"

	"gangsaur.com/billing-exercise/internal/db/psql"
)

type LoanStore interface {
	GetLoan(ctx context.Context, id int) (psql.Loan, error)
}

type LoanService struct {
	store LoanStore
}

func NewLoanService(l LoanStore) *LoanService {
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
