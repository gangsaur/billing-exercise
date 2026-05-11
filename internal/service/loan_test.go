package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"gangsaur.com/billing-exercise/internal/repository/db/psql"
	"gangsaur.com/billing-exercise/internal/service"
	storeMocks "gangsaur.com/billing-exercise/internal/service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoanService_GetLoan(t *testing.T) {
	sampleTime := time.Now()

	tests := []struct {
		name                string
		mockStore           *storeMocks.MockStore
		id                  int
		mockGetLoanResponse psql.Loan
		mockGetLoanErr      error
		want                psql.Loan
		wantErr             bool
	}{
		{
			name:      "success case",
			id:        1,
			mockStore: storeMocks.NewMockStore(t),
			mockGetLoanResponse: psql.Loan{
				Id:                1,
				Duration:          50,
				PrincipalAmount:   5000000,
				OutstandingAmount: 5000000,
				Status:            0,
				InterestRate:      10,
				UserId:            1,
				CreatedAt:         sampleTime,
				UpdatedAt:         sampleTime,
			},
			want: psql.Loan{
				Id:                1,
				Duration:          50,
				PrincipalAmount:   5000000,
				OutstandingAmount: 5000000,
				Status:            0,
				InterestRate:      10,
				UserId:            1,
				CreatedAt:         sampleTime,
				UpdatedAt:         sampleTime,
			},
		},
		{
			name:           "error case, GetLoan throw error",
			id:             2,
			mockStore:      storeMocks.NewMockStore(t),
			mockGetLoanErr: errors.New("GetLoan error"),
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := service.NewLoanService(tt.mockStore)

			tt.mockStore.On("GetLoan", mock.Anything, tt.id).Return(tt.mockGetLoanResponse, tt.mockGetLoanErr)

			got, gotErr := l.GetLoan(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestLoanService_PayLoan(t *testing.T) {
	type loanAndLoanPayments struct {
		Loan         psql.Loan
		LoanPayments []psql.LoanPayment
	}

	tests := []struct {
		name                                              string
		mockStore                                         *storeMocks.MockStore
		id                                                int
		amount                                            int
		mockGetLoanAndLoanPaymentsByStatusDueDateResponse loanAndLoanPayments
		mockGetLoanAndLoanPaymentsByStatusDueDateErr      error
		flagInvalidAmount                                 bool
		mockBeginErr                                      error
		paramsOutstandingDeduction                        int
		paramsOutstandingAmount                           int
		mockReduceLoanOutstandingAmountTxErr              error
		paramsLoanPaymentIds                              []int
		mockUpdateLoanPaymentStatusPaidTxErr              error
		mockCommitErr                                     error
		mockGetLoanResponse                               psql.Loan
		mockGetLoanErr                                    error
		want                                              psql.Loan
		wantErr                                           bool
	}{
		{
			name:      "success case, single payment",
			mockStore: storeMocks.NewMockStore(t),
			id:        1,
			amount:    110000,
			mockGetLoanAndLoanPaymentsByStatusDueDateResponse: loanAndLoanPayments{
				Loan:         psql.Loan{Id: 1, InterestRate: 10, OutstandingAmount: 3900000},
				LoanPayments: []psql.LoanPayment{{Id: 1, LoanId: 1, Amount: 110000}},
			},
			paramsOutstandingDeduction: 100000,
			paramsOutstandingAmount:    3900000,
			paramsLoanPaymentIds:       []int{1},
			mockGetLoanResponse:        psql.Loan{Id: 1, InterestRate: 10, OutstandingAmount: 3800000},
			want:                       psql.Loan{Id: 1, InterestRate: 10, OutstandingAmount: 3800000},
		},
		{
			name:      "success case, multiple payments",
			mockStore: storeMocks.NewMockStore(t),
			id:        1,
			amount:    330000,
			mockGetLoanAndLoanPaymentsByStatusDueDateResponse: loanAndLoanPayments{
				Loan: psql.Loan{Id: 1, InterestRate: 10, OutstandingAmount: 300000},
				LoanPayments: []psql.LoanPayment{
					{Id: 1, LoanId: 1, Amount: 110000},
					{Id: 99, LoanId: 1, Amount: 110000},
					{Id: 2, LoanId: 1, Amount: 110000},
				},
			},
			paramsOutstandingDeduction: 300000,
			paramsOutstandingAmount:    300000,
			paramsLoanPaymentIds:       []int{1, 99, 2},
			mockGetLoanResponse:        psql.Loan{Id: 1, InterestRate: 10, OutstandingAmount: 3800000},
			want:                       psql.Loan{Id: 1, InterestRate: 10, OutstandingAmount: 3800000},
		},
		{
			name:      "error case, GetLoanAndLoanPaymentsByStatusDueDate error",
			mockStore: storeMocks.NewMockStore(t),
			id:        1,
			amount:    220000,
			mockGetLoanAndLoanPaymentsByStatusDueDateErr: errors.New("GetLoanAndLoanPaymentsByStatusDueDate error"),
			wantErr: true,
		},
		{
			name:      "error case, invalid amount",
			mockStore: storeMocks.NewMockStore(t),
			id:        1,
			amount:    550000,
			mockGetLoanAndLoanPaymentsByStatusDueDateResponse: loanAndLoanPayments{
				Loan: psql.Loan{Id: 1, InterestRate: 10, OutstandingAmount: 300000},
				LoanPayments: []psql.LoanPayment{
					{Id: 1000, LoanId: 1, Amount: 110000},
					{Id: 99, LoanId: 1, Amount: 110000},
				},
			},
			flagInvalidAmount: true,
			wantErr:           true,
		},
		{
			name:      "error case, Begin error",
			mockStore: storeMocks.NewMockStore(t),
			id:        1,
			amount:    220000,
			mockGetLoanAndLoanPaymentsByStatusDueDateResponse: loanAndLoanPayments{
				Loan: psql.Loan{Id: 1, InterestRate: 10, OutstandingAmount: 300000},
				LoanPayments: []psql.LoanPayment{
					{Id: 1000, LoanId: 1, Amount: 110000},
					{Id: 99, LoanId: 1, Amount: 110000},
				},
			},
			mockBeginErr: errors.New("Begin error"),
			wantErr:      true,
		},
		{
			name:      "error case, ReduceLoanOutstandingAmountTx error",
			mockStore: storeMocks.NewMockStore(t),
			id:        1,
			amount:    220000,
			mockGetLoanAndLoanPaymentsByStatusDueDateResponse: loanAndLoanPayments{
				Loan: psql.Loan{Id: 1, InterestRate: 10, OutstandingAmount: 300000},
				LoanPayments: []psql.LoanPayment{
					{Id: 1000, LoanId: 1, Amount: 110000},
					{Id: 99, LoanId: 1, Amount: 110000},
				},
			},
			paramsOutstandingDeduction:           200000,
			paramsOutstandingAmount:              300000,
			mockReduceLoanOutstandingAmountTxErr: errors.New("ReduceLoanOutstandingAmountTx error"),
			wantErr:                              true,
		},
		{
			name:      "error case, ReduceLoanOutstandingAmountStatusPaidTx error",
			mockStore: storeMocks.NewMockStore(t),
			id:        1,
			amount:    220000,
			mockGetLoanAndLoanPaymentsByStatusDueDateResponse: loanAndLoanPayments{
				Loan: psql.Loan{Id: 1, InterestRate: 10, OutstandingAmount: 200000},
				LoanPayments: []psql.LoanPayment{
					{Id: 1000, LoanId: 1, Amount: 110000},
					{Id: 99, LoanId: 1, Amount: 110000},
				},
			},
			paramsOutstandingDeduction:           200000,
			paramsOutstandingAmount:              200000,
			mockReduceLoanOutstandingAmountTxErr: errors.New("ReduceLoanOutstandingAmountStatusPaidTx error"),
			wantErr:                              true,
		},
		{
			name:      "error case, UpdateLoanPaymentStatusPaidTx error",
			mockStore: storeMocks.NewMockStore(t),
			id:        1,
			amount:    220000,
			mockGetLoanAndLoanPaymentsByStatusDueDateResponse: loanAndLoanPayments{
				Loan: psql.Loan{Id: 1, InterestRate: 10, OutstandingAmount: 300000},
				LoanPayments: []psql.LoanPayment{
					{Id: 1000, LoanId: 1, Amount: 110000},
					{Id: 99, LoanId: 1, Amount: 110000},
				},
			},
			paramsOutstandingDeduction:           200000,
			paramsOutstandingAmount:              300000,
			paramsLoanPaymentIds:                 []int{1000, 99},
			mockUpdateLoanPaymentStatusPaidTxErr: errors.New("UpdateLoanPaymentStatusPaidTx error"),
			wantErr:                              true,
		},
		{
			name:      "error case, Commit error",
			mockStore: storeMocks.NewMockStore(t),
			id:        1,
			amount:    220000,
			mockGetLoanAndLoanPaymentsByStatusDueDateResponse: loanAndLoanPayments{
				Loan: psql.Loan{Id: 1, InterestRate: 10, OutstandingAmount: 300000},
				LoanPayments: []psql.LoanPayment{
					{Id: 1000, LoanId: 1, Amount: 110000},
					{Id: 99, LoanId: 1, Amount: 110000},
				},
			},
			paramsOutstandingDeduction: 200000,
			paramsOutstandingAmount:    300000,
			paramsLoanPaymentIds:       []int{1000, 99},
			mockCommitErr:              errors.New("Commit error"),
			wantErr:                    true,
		},
		{
			name:      "error case, GetLoan error",
			mockStore: storeMocks.NewMockStore(t),
			id:        1,
			amount:    220000,
			mockGetLoanAndLoanPaymentsByStatusDueDateResponse: loanAndLoanPayments{
				Loan: psql.Loan{Id: 1, InterestRate: 10, OutstandingAmount: 300000},
				LoanPayments: []psql.LoanPayment{
					{Id: 1000, LoanId: 1, Amount: 110000},
					{Id: 99, LoanId: 1, Amount: 110000},
				},
			},
			paramsOutstandingDeduction: 200000,
			paramsOutstandingAmount:    300000,
			paramsLoanPaymentIds:       []int{1000, 99},
			mockGetLoanErr:             errors.New("GetLoan error"),
			wantErr:                    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := service.NewLoanService(tt.mockStore)

			// Mocking
			tt.mockStore.On("GetLoanAndLoanPaymentsByStatusDueDate", mock.Anything, tt.id, psql.LoanPaymentStatusScheduled, mock.Anything, true).Return(
				tt.mockGetLoanAndLoanPaymentsByStatusDueDateResponse.Loan,
				tt.mockGetLoanAndLoanPaymentsByStatusDueDateResponse.LoanPayments,
				tt.mockGetLoanAndLoanPaymentsByStatusDueDateErr)

			if tt.mockGetLoanAndLoanPaymentsByStatusDueDateErr == nil && !tt.flagInvalidAmount {
				tt.mockStore.On("Begin", mock.Anything).Return(nil, tt.mockBeginErr)

				if tt.mockBeginErr == nil {
					tt.mockStore.On("Rollback", mock.Anything, mock.Anything).Return(nil)

					if tt.paramsOutstandingDeduction == tt.paramsOutstandingAmount {
						tt.mockStore.On("ReduceLoanOutstandingAmountStatusPaidTx", mock.Anything, mock.Anything, tt.paramsOutstandingDeduction, tt.id, tt.paramsOutstandingAmount).
							Return(tt.mockReduceLoanOutstandingAmountTxErr)
					} else {
						tt.mockStore.On("ReduceLoanOutstandingAmountTx", mock.Anything, mock.Anything, tt.paramsOutstandingDeduction, tt.id, tt.paramsOutstandingAmount).
							Return(tt.mockReduceLoanOutstandingAmountTxErr)
					}

					if tt.mockReduceLoanOutstandingAmountTxErr == nil {
						tt.mockStore.On("UpdateLoanPaymentStatusPaidTx", mock.Anything, mock.Anything, tt.paramsLoanPaymentIds).
							Return(tt.mockUpdateLoanPaymentStatusPaidTxErr)

						if tt.mockUpdateLoanPaymentStatusPaidTxErr == nil {
							tt.mockStore.On("Commit", mock.Anything, mock.Anything).Return(tt.mockCommitErr)

							if tt.mockCommitErr == nil {
								tt.mockStore.On("GetLoan", mock.Anything, tt.id).Return(tt.mockGetLoanResponse, tt.mockGetLoanErr)
							}
						}
					}
				}
			}

			//Call
			got, gotErr := l.PayLoan(context.Background(), tt.id, tt.amount)
			if tt.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
