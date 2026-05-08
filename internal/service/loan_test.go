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
			mockGetLoanErr: nil,
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
			wantErr: false,
		},
		{
			name:                "error case, GetLoan throw error",
			id:                  2,
			mockStore:           storeMocks.NewMockStore(t),
			mockGetLoanResponse: psql.Loan{},
			mockGetLoanErr:      errors.New("GetLoan error"),
			want:                psql.Loan{},
			wantErr:             true,
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
