package service

import "context"

type LoanData struct {
}

type LoanService struct {
}

func (l *LoanService) GetLoan(ctx context.Context, id int) (LoanData, error) {
	return nil, nil
}
