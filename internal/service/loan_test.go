package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"gangsaur.com/billing-exercise/internal/repository/db/psql"
	storeMocks "gangsaur.com/billing-exercise/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoanService_GetLoan(t *testing.T) {
	mockStore := storeMocks.NewMockStore(t) // Use shared mocks
	sampleTime := time.Now()

	tests := []struct {
		name                string
		id                  int
		mockGetLoanResponse psql.Loan
		mockGetLoanErr      error
		want                psql.Loan
		wantErr             bool
	}{
		{
			name: "success case",
			id:   1,
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
			mockGetLoanResponse: psql.Loan{},
			mockGetLoanErr:      errors.New("sample-error"),
			want:                psql.Loan{},
			wantErr:             true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLoanService(mockStore)

			mockStore.On("GetLoan", mock.Anything, tt.id).Return(tt.mockGetLoanResponse, tt.mockGetLoanErr)

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
