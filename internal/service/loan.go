package service

import (
	"context"
	"fmt"
	"slices"
	"time"

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

func (l *LoanService) PayLoan(ctx context.Context, id int, amount int) (psql.Loan, error) { // Refactor and fix later after test is set up
	// GetLoanData
	loanData, err := l.store.GetLoan(ctx, id)
	if err != nil {
		return psql.Loan{}, err
	}

	// Get unpaid loan payments and calculate the needed payment
	unpaidLoanPayments, err := l.store.GetLoanPaymentsByLoanIdsStatusDueDate(ctx, []int{id}, psql.LoanPaymentStatusScheduled, time.Now().AddDate(0, 0, 7), true)
	if err != nil {
		return psql.Loan{}, err
	}

	totalUnpaid := 0
	loanPaymentIds := make([]int, 0, len(unpaidLoanPayments))
	for _, upl := range unpaidLoanPayments {
		totalUnpaid += upl.Amount
		loanPaymentIds = append(loanPaymentIds, upl.Id)
	}
	if amount != totalUnpaid {
		return psql.Loan{}, fmt.Errorf("Invalid payment amount, payment amount must be %d", totalUnpaid)
	}
	slices.Sort(loanPaymentIds)

	// Make Payment - Update to use transaction laters
	outstandingDeduction := float32(amount) / (1.0 + (loanData.InterestRate / 100.0))
	err = l.store.PayLoan(ctx, id, loanPaymentIds, int(outstandingDeduction), loanData.OutstandingAmount)
	if err != nil {
		return psql.Loan{}, err
	}

	// Fetch latest loan data
	loanData, err = l.store.GetLoan(ctx, id)
	if err != nil {
		return psql.Loan{}, err
	}

	return loanData, nil
}

func (l *LoanService) PayLoanV2(ctx context.Context, id int, amount int) (psql.Loan, error) {
	// Need centralized lock or payload needs to specify loanPaymentIds to handle race condition

	// Get Loan and their LoanPayments
	loan, unpaidLoanPayments, err := l.store.GetLoanAndLoanPaymentsByStatusDueDate(ctx, id, psql.LoanPaymentStatusScheduled, time.Now().AddDate(0, 0, 7), true)
	if err != nil {
		return psql.Loan{}, err
	}

	// Calculate Unpaid
	totalUnpaid := 0
	loanPaymentIds := make([]int, 0, len(unpaidLoanPayments))
	for _, upl := range unpaidLoanPayments {
		totalUnpaid += upl.Amount
		loanPaymentIds = append(loanPaymentIds, upl.Id)
	}
	if amount != totalUnpaid {
		return psql.Loan{}, fmt.Errorf("Invalid payment amount, payment amount must be %d", totalUnpaid)
	}
	slices.Sort(loanPaymentIds)

	// Make Payment
	tx, err := l.store.Begin(ctx)
	if err != nil {
		return psql.Loan{}, err
	}
	defer l.store.Rollback(ctx, tx)

	// Update Loan
	outstandingDeduction := int(float32(amount) / (1.0 + (loan.InterestRate / 100.0)))
	if outstandingDeduction == loan.OutstandingAmount {
		err = l.store.ReduceLoanOutstandingAmountStatusPaidTx(ctx, tx, outstandingDeduction, id, loan.OutstandingAmount)
	} else {
		err = l.store.ReduceLoanOutstandingAmountTx(ctx, tx, outstandingDeduction, id, loan.OutstandingAmount)
	}
	if err != nil {
		return psql.Loan{}, err
	}

	// Update Loan Payments
	err = l.store.UpdateLoanPaymentStatusPaidTx(ctx, tx, loanPaymentIds)
	if err != nil {
		return psql.Loan{}, err
	}

	// Commit
	err = l.store.Commit(ctx, tx)
	if err != nil {
		return psql.Loan{}, err
	}

	// Fetch latest loan data
	loan, err = l.store.GetLoan(ctx, id)
	if err != nil {
		return psql.Loan{}, err
	}

	return loan, nil
}
