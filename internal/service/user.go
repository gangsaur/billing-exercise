package service

import (
	"context"
	"time"

	"gangsaur.com/billing-exercise/internal/db/psql"
)

const DelinquentThreshold = 2

type UserService struct {
	store Store
}

func NewUserService(u Store) *UserService {
	return &UserService{
		store: u,
	}
}

func (u *UserService) GetUser(ctx context.Context, id int) (psql.User, bool, error) {
	// Get user data
	userData, err := u.store.GetUser(ctx, id)
	if err != nil {
		return psql.User{}, false, err
	}

	// Get user's open loan
	loans, err := u.store.GetLoanByUserIdAndStatus(ctx, id, psql.LoanStatusOpen)
	if err != nil {
		return psql.User{}, false, err
	}

	loanIds := make([]int, 0, len(loans))
	loanUnpaidPaymentCount := make(map[int]int)
	for _, l := range loans {
		loanIds = append(loanIds, l.Id)
		loanUnpaidPaymentCount[l.Id] = 0
	}

	// Get user's unpaid loan payments
	loanPayments, err := u.store.GetLoanPaymentsByLoanIdsStatusDueDate(ctx, loanIds, psql.LoanPaymentStatusScheduled, time.Now().AddDate(0, 0, 7), true)
	if err != nil {
		return psql.User{}, false, err
	}

	delinquentStatus := false
	for _, lp := range loanPayments {
		loanUnpaidPaymentCount[lp.LoanId]++
		if loanUnpaidPaymentCount[lp.LoanId] > DelinquentThreshold {
			delinquentStatus = true
			break
		}
	}

	return userData, delinquentStatus, nil
}
