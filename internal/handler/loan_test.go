package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"gangsaur.com/billing-exercise/internal/handler"
	"gangsaur.com/billing-exercise/internal/repository/db/psql"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoanHandler_GetLoan(t *testing.T) {
	sampleTime := time.Now()

	tests := []struct {
		name                string
		mockService         *handler.MockLoanService
		requestParamsId     int
		mockGetLoanResponse psql.Loan
		mockGetLoanErr      error
		wantStatus          int
		wantBodyChecker     func(*testing.T, []byte)
	}{
		{
			name:            "success case",
			mockService:     handler.NewMockLoanService(t),
			requestParamsId: 1,
			mockGetLoanResponse: psql.Loan{
				Id:                1,
				Duration:          50,
				PrincipalAmount:   5000000,
				OutstandingAmount: 200000,
				Status:            0,
				InterestRate:      10.0,
				UserId:            1,
				CreatedAt:         sampleTime,
				UpdatedAt:         sampleTime,
			},
			wantStatus: 200,
			wantBodyChecker: func(t *testing.T, body []byte) {
				var loanResponse handler.LoanResponse
				_ = json.Unmarshal(body, &loanResponse)

				assert.Equal(t, 1, loanResponse.Id)
				assert.Equal(t, 50, loanResponse.Duration)
				assert.Equal(t, 5000000, loanResponse.PrincipalAmount)
				assert.Equal(t, 200000, loanResponse.OutstandingAmount)
				assert.Equal(t, 0, loanResponse.Status)
				assert.Equal(t, float32(10), loanResponse.InterestRate)
				assert.Equal(t, 1, loanResponse.UserId)
			},
		},
		{
			name:            "error case, invalid id",
			mockService:     handler.NewMockLoanService(t),
			requestParamsId: -1,
			wantStatus:      500,
			wantBodyChecker: func(t *testing.T, body []byte) {},
		},
		{
			name:            "error case, GetLoan error",
			mockService:     handler.NewMockLoanService(t),
			requestParamsId: 1,
			mockGetLoanErr:  errors.New("GetLoan error"),
			wantStatus:      500,
			wantBodyChecker: func(t *testing.T, body []byte) {},
		},
		{
			name:            "error case, GetLoan not found",
			mockService:     handler.NewMockLoanService(t),
			requestParamsId: 1,
			mockGetLoanErr:  psql.ErrNotFound,
			wantStatus:      404,
			wantBodyChecker: func(t *testing.T, body []byte) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := handler.NewLoanHandler(tt.mockService)

			// Construct test request and setup mocks
			requestParamsIdStr := strconv.Itoa(tt.requestParamsId)
			r := httptest.NewRequest("GET", "/loan/"+requestParamsIdStr, nil)

			if tt.requestParamsId != -1 {
				r.SetPathValue("id", requestParamsIdStr)
				tt.mockService.On("GetLoan", mock.Anything, tt.requestParamsId).
					Return(tt.mockGetLoanResponse, tt.mockGetLoanErr)
			} else {
				r.SetPathValue("id", "invalid id")
			}

			w := httptest.NewRecorder()

			// Call the handlerFunc
			l.GetLoan()(w, r)

			// Check the result
			res := w.Result()
			body, _ := io.ReadAll(res.Body)
			res.Body.Close()

			assert.Equal(t, tt.wantStatus, res.StatusCode)
			tt.wantBodyChecker(t, bytes.TrimSpace(body))
		})
	}
}
