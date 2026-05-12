package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"gangsaur.com/billing-exercise/internal/repository/db/psql"
	"gangsaur.com/billing-exercise/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_GetUser(t *testing.T) {
	sampleTime := time.Now()

	tests := []struct {
		name                                              string
		mockStore                                         *service.MockStore
		id                                                int
		mockGetUserResponse                               psql.User
		mockGetUserErr                                    error
		mockGetLoanByUserIdAndStatusResponse              []psql.Loan
		mockGetLoanByUserIdAndStatusErr                   error
		paramsLoandIds                                    []int
		mockGetLoanPaymentsByLoanIdsStatusDueDateResponse []psql.LoanPayment
		mockGetLoanPaymentsByLoanIdsStatusDueDateErr      error
		wantUser                                          psql.User
		wantDelinquentStatus                              bool
		wantErr                                           bool
	}{
		{
			name:      "success case, not delinquent",
			mockStore: service.NewMockStore(t),
			id:        1,
			mockGetUserResponse: psql.User{
				Id:        1,
				CreatedAt: sampleTime,
				UpdatedAt: sampleTime,
			},
			mockGetLoanByUserIdAndStatusResponse: []psql.Loan{
				{Id: 1},
				{Id: 2},
			},
			paramsLoandIds: []int{1, 2},
			mockGetLoanPaymentsByLoanIdsStatusDueDateResponse: []psql.LoanPayment{
				{Id: 1, LoanId: 1},
				{Id: 2, LoanId: 1},
				{Id: 3, LoanId: 2},
			},
			wantUser: psql.User{
				Id:        1,
				CreatedAt: sampleTime,
				UpdatedAt: sampleTime,
			},
			wantDelinquentStatus: false,
			wantErr:              false,
		},
		{
			name:      "success case, delinquent",
			mockStore: service.NewMockStore(t),
			id:        1,
			mockGetUserResponse: psql.User{
				Id:        1,
				CreatedAt: sampleTime,
				UpdatedAt: sampleTime,
			},
			mockGetLoanByUserIdAndStatusResponse: []psql.Loan{
				{Id: 99},
				{Id: 1000},
			},
			paramsLoandIds: []int{99, 1000},
			mockGetLoanPaymentsByLoanIdsStatusDueDateResponse: []psql.LoanPayment{
				{Id: 1, LoanId: 1},
				{Id: 2, LoanId: 1},
				{Id: 3, LoanId: 1},
			},
			wantUser: psql.User{
				Id:        1,
				CreatedAt: sampleTime,
				UpdatedAt: sampleTime,
			},
			wantDelinquentStatus: true,
			wantErr:              false,
		},
		{
			name:           "error case, GetUser error",
			mockStore:      service.NewMockStore(t),
			id:             1,
			mockGetUserErr: errors.New("GetUser error"),
			wantErr:        true,
		},
		{
			name:      "error case, GetLoanByUserIdAndStatus error",
			mockStore: service.NewMockStore(t),
			id:        1,
			mockGetUserResponse: psql.User{
				Id:        1,
				CreatedAt: sampleTime,
				UpdatedAt: sampleTime,
			},
			mockGetLoanByUserIdAndStatusErr: errors.New("GetLoanByUserIdAndStatus error"),
			wantErr:                         true,
		},
		{
			name:      "error case, GetLoanByUserIdAndStatus error",
			mockStore: service.NewMockStore(t),
			id:        1,
			mockGetUserResponse: psql.User{
				Id:        1,
				CreatedAt: sampleTime,
				UpdatedAt: sampleTime,
			},
			mockGetLoanByUserIdAndStatusResponse: []psql.Loan{
				{Id: 77},
				{Id: 3133},
			},
			paramsLoandIds: []int{77, 3133},
			mockGetLoanPaymentsByLoanIdsStatusDueDateErr: errors.New("GetLoanByUserIdAndStatus error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := service.NewUserService(tt.mockStore)

			tt.mockStore.On("GetUser", mock.Anything, tt.id).Return(tt.mockGetUserResponse, tt.mockGetUserErr)
			if tt.mockGetUserErr == nil {
				tt.mockStore.On("GetLoanByUserIdAndStatus", mock.Anything, tt.id, psql.LoanStatusOpen).
					Return(tt.mockGetLoanByUserIdAndStatusResponse, tt.mockGetLoanByUserIdAndStatusErr)
			}
			if tt.mockGetUserErr == nil && tt.mockGetLoanByUserIdAndStatusErr == nil {
				tt.mockStore.On("GetLoanPaymentsByLoanIdsStatusDueDate", mock.Anything, tt.paramsLoandIds, psql.LoanPaymentStatusScheduled, mock.Anything, true).
					Return(tt.mockGetLoanPaymentsByLoanIdsStatusDueDateResponse, tt.mockGetLoanPaymentsByLoanIdsStatusDueDateErr)
			}

			gotUser, gotDelinquenStatus, gotErr := u.GetUser(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
				assert.Equal(t, tt.wantUser, gotUser)
				assert.Equal(t, tt.wantDelinquentStatus, gotDelinquenStatus)
			}
		})
	}
}
