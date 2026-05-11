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
	tests := []struct {
		name                                              string
		mockStore                                         *storeMocks.MockStore
		id                                                int
		amount                                            int
		mockGetLoanResponse                               []psql.Loan
		mockGetLoanErr                                    []error
		mockGetLoanPaymentsByLoanIdsStatusDueDateResponse []psql.LoanPayment
		mockGetLoanPaymentsByLoanIdsStatusDueDateErr      error
		flagInvalidAmount                                 bool
		paramsLoanPaymentIds                              []int
		paramsOutstandingDeduction                        int
		paramsOutstandingAmount                           int
		mockPayLoanErr                                    error
		want                                              psql.Loan
		wantErr                                           bool
	}{
		{
			name:      "success case, single payment",
			mockStore: storeMocks.NewMockStore(t),
			id:        1,
			amount:    110000,
			mockGetLoanResponse: []psql.Loan{
				{
					Id:                1,
					InterestRate:      10,
					OutstandingAmount: 3900000,
				},
				{
					Id:                1,
					InterestRate:      10,
					OutstandingAmount: 3800000,
				},
			},
			mockGetLoanErr: []error{nil, nil},
			mockGetLoanPaymentsByLoanIdsStatusDueDateResponse: []psql.LoanPayment{
				{Id: 1, LoanId: 1, Amount: 110000},
			},
			paramsLoanPaymentIds:       []int{1},
			paramsOutstandingDeduction: 100000,
			paramsOutstandingAmount:    3900000,
			want: psql.Loan{
				Id:                1,
				InterestRate:      10,
				OutstandingAmount: 3800000,
			},
		},
		{
			name:      "success case, multiple payment",
			mockStore: storeMocks.NewMockStore(t),
			id:        1,
			amount:    315000,
			mockGetLoanResponse: []psql.Loan{
				{
					Id:                1,
					InterestRate:      5,
					OutstandingAmount: 3800000,
				},
				{
					Id:                1,
					InterestRate:      5,
					OutstandingAmount: 3500000,
				},
			},
			mockGetLoanErr: []error{nil, nil},
			mockGetLoanPaymentsByLoanIdsStatusDueDateResponse: []psql.LoanPayment{
				{Id: 33, LoanId: 1, Amount: 105000},
				{Id: 34, LoanId: 1, Amount: 105000},
				{Id: 35, LoanId: 1, Amount: 105000},
			},
			paramsLoanPaymentIds:       []int{33, 34, 35},
			paramsOutstandingDeduction: 300000,
			paramsOutstandingAmount:    3800000,
			want: psql.Loan{
				Id:                1,
				InterestRate:      5,
				OutstandingAmount: 3500000,
			},
		},
		{
			name:      "error case, 1st GetLoan error",
			mockStore: storeMocks.NewMockStore(t),
			id:        1,
			amount:    1000,
			mockGetLoanResponse: []psql.Loan{
				{},
				{},
			},
			mockGetLoanErr: []error{errors.New("GetLoan error"), nil},
			wantErr:        true,
		},
		{
			name:      "error case, GetLoanPaymentsByLoanIdsStatusDueDate error",
			mockStore: storeMocks.NewMockStore(t),
			id:        1,
			amount:    200000,
			mockGetLoanResponse: []psql.Loan{
				{
					Id:                1,
					InterestRate:      0,
					OutstandingAmount: 3500000,
				},
				{},
			},
			mockGetLoanErr: []error{nil, nil},
			mockGetLoanPaymentsByLoanIdsStatusDueDateErr: errors.New("GetLoanPaymentsByLoanIdsStatusDueDate error"),
			wantErr: true,
		},
		{
			name:      "error case, invalid payment amount",
			mockStore: storeMocks.NewMockStore(t),
			id:        1,
			amount:    500000,
			mockGetLoanResponse: []psql.Loan{
				{
					Id:                1,
					InterestRate:      0,
					OutstandingAmount: 3500000,
				},
				{},
			},
			mockGetLoanErr: []error{nil, nil},
			mockGetLoanPaymentsByLoanIdsStatusDueDateResponse: []psql.LoanPayment{
				{Id: 35, LoanId: 1, Amount: 100000},
				{Id: 36, LoanId: 1, Amount: 100000},
			},
			flagInvalidAmount: true,
			wantErr:           true,
		},
		{
			name:      "error case, PayLoan error",
			mockStore: storeMocks.NewMockStore(t),
			id:        1,
			amount:    200000,
			mockGetLoanResponse: []psql.Loan{
				{
					Id:                1,
					InterestRate:      0,
					OutstandingAmount: 3800000,
				},
				{},
			},
			mockGetLoanErr: []error{nil, nil},
			mockGetLoanPaymentsByLoanIdsStatusDueDateResponse: []psql.LoanPayment{
				{Id: 35, LoanId: 1, Amount: 100000},
				{Id: 36, LoanId: 1, Amount: 100000},
			},
			paramsLoanPaymentIds:       []int{35, 36},
			paramsOutstandingDeduction: 200000,
			paramsOutstandingAmount:    3800000,
			mockPayLoanErr:             errors.New("PayLoan error"),
			wantErr:                    true,
		},
		{
			name:      "error case, 2nd GetLoan error",
			mockStore: storeMocks.NewMockStore(t),
			id:        1,
			amount:    200000,
			mockGetLoanResponse: []psql.Loan{
				{
					Id:                1,
					InterestRate:      0,
					OutstandingAmount: 3800000,
				},
				{},
			},
			mockGetLoanErr: []error{nil, errors.New("GetLoan error")},
			mockGetLoanPaymentsByLoanIdsStatusDueDateResponse: []psql.LoanPayment{
				{Id: 35, LoanId: 1, Amount: 100000},
				{Id: 36, LoanId: 1, Amount: 100000},
			},
			paramsLoanPaymentIds:       []int{35, 36},
			paramsOutstandingDeduction: 200000,
			paramsOutstandingAmount:    3800000,
			wantErr:                    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := service.NewLoanService(tt.mockStore)

			tt.mockStore.On("GetLoan", mock.Anything, tt.id).Return(tt.mockGetLoanResponse[0], tt.mockGetLoanErr[0]).Once()
			if tt.mockGetLoanErr[0] == nil {
				tt.mockStore.On("GetLoanPaymentsByLoanIdsStatusDueDate", mock.Anything, []int{tt.id}, psql.LoanPaymentStatusScheduled, mock.Anything, true).
					Return(tt.mockGetLoanPaymentsByLoanIdsStatusDueDateResponse, tt.mockGetLoanPaymentsByLoanIdsStatusDueDateErr)
			}
			if tt.mockGetLoanErr[0] == nil && tt.mockGetLoanPaymentsByLoanIdsStatusDueDateErr == nil && !tt.flagInvalidAmount {
				tt.mockStore.On("PayLoan", mock.Anything, tt.id, tt.paramsLoanPaymentIds, tt.paramsOutstandingDeduction, tt.paramsOutstandingAmount).
					Return(tt.mockPayLoanErr)
			}
			if tt.mockGetLoanErr[0] == nil && tt.mockGetLoanPaymentsByLoanIdsStatusDueDateErr == nil && !tt.flagInvalidAmount && tt.mockPayLoanErr == nil {
				tt.mockStore.On("GetLoan", mock.Anything, tt.id).Return(tt.mockGetLoanResponse[1], tt.mockGetLoanErr[1]).Once()
			}

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

func TestLoanService_PayLoanV2(t *testing.T) {
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
				tt.mockStore.On("Rollback", mock.Anything, mock.Anything).Return(nil)

				if tt.mockBeginErr == nil {
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
							tt.mockStore.On("Commit", mock.Anything, mock.Anything).Return(tt.mockBeginErr)

							if tt.mockBeginErr == nil {
								tt.mockStore.On("GetLoan", mock.Anything, tt.id).Return(tt.mockGetLoanResponse, tt.mockGetLoanErr)
							}
						}
					}
				}
			}

			//Call
			got, gotErr := l.PayLoanV2(context.Background(), tt.id, tt.amount)
			if tt.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
