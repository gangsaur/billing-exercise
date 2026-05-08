package psql

import (
	"context"
	"time"
)

type Psql struct {
}

type Loan struct {
	Id                int
	Duration          int
	PrincipalAmount   int
	OutstandingAmount int
	InterestRate      float32
	UserId            int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (p *Psql) GetLoan(ctx context.Context, id int) (Loan, error) {
	return Loan{}, nil
}
