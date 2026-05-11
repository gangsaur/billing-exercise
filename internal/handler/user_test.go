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
	mock "github.com/stretchr/testify/mock"
)

func TestUserHandler_GetUser(t *testing.T) {
	type mockGetUserResponse struct {
		user       psql.User
		delinquent bool
	}

	sampleTime := time.Now()

	tests := []struct {
		name                string
		mockService         *handler.MockUserService
		requestParamsId     int
		mockGetUserResponse mockGetUserResponse
		mockGetUserErr      error
		wantStatus          int
		wantBodyChecker     func(*testing.T, []byte)
	}{
		{
			name:            "success case, not delinquent",
			mockService:     handler.NewMockUserService(t),
			requestParamsId: 1,
			mockGetUserResponse: mockGetUserResponse{
				user:       psql.User{Id: 1, CreatedAt: sampleTime, UpdatedAt: sampleTime},
				delinquent: false,
			},
			wantStatus: 200,
			wantBodyChecker: func(t *testing.T, body []byte) {
				var userResponseDelinquent handler.UserResponseDelinquent
				_ = json.Unmarshal(body, &userResponseDelinquent)

				assert.Equal(t, 1, userResponseDelinquent.Id)
				assert.Equal(t, false, userResponseDelinquent.Delinquent)
			},
		},
		{
			name:            "success case, is delinquent",
			mockService:     handler.NewMockUserService(t),
			requestParamsId: 2,
			mockGetUserResponse: mockGetUserResponse{
				user:       psql.User{Id: 2, CreatedAt: sampleTime, UpdatedAt: sampleTime},
				delinquent: true,
			},
			wantStatus: 200,
			wantBodyChecker: func(t *testing.T, body []byte) {
				var userResponseDelinquent handler.UserResponseDelinquent
				_ = json.Unmarshal(body, &userResponseDelinquent)

				assert.Equal(t, 2, userResponseDelinquent.Id)
				assert.Equal(t, true, userResponseDelinquent.Delinquent)
			},
		},
		{
			name:            "error case, invalid id",
			mockService:     handler.NewMockUserService(t),
			requestParamsId: -1,
			wantStatus:      500,
			wantBodyChecker: func(t *testing.T, body []byte) {},
		},
		{
			name:            "error case, GetUser error",
			mockService:     handler.NewMockUserService(t),
			requestParamsId: 1,
			mockGetUserErr:  errors.New("GetUser error"),
			wantStatus:      500,
			wantBodyChecker: func(t *testing.T, body []byte) {},
		},
		{
			name:            "error case, GetUser not found",
			mockService:     handler.NewMockUserService(t),
			requestParamsId: 1,
			mockGetUserErr:  psql.ErrNotFound,
			wantStatus:      404,
			wantBodyChecker: func(t *testing.T, body []byte) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := handler.NewUserHandler(tt.mockService)

			// Construct test request and setup mocks
			requestParamsIdStr := strconv.Itoa(tt.requestParamsId)
			r := httptest.NewRequest("GET", "/user/"+requestParamsIdStr, nil)

			if tt.requestParamsId != -1 {
				r.SetPathValue("id", requestParamsIdStr)
				tt.mockService.On("GetUser", mock.Anything, tt.requestParamsId).
					Return(tt.mockGetUserResponse.user, tt.mockGetUserResponse.delinquent, tt.mockGetUserErr)
			} else {
				r.SetPathValue("id", "invalid id")
			}

			w := httptest.NewRecorder()

			// Call the handlerFunc
			u.GetUser()(w, r)

			// Check the result
			res := w.Result()
			body, _ := io.ReadAll(res.Body)
			res.Body.Close()

			assert.Equal(t, tt.wantStatus, res.StatusCode)
			tt.wantBodyChecker(t, bytes.TrimSpace(body))
		})
	}
}
