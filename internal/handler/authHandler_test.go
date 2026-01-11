package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/jacobpq/soccer-manager/internal/api"
	"github.com/jacobpq/soccer-manager/internal/domain/models"
	"github.com/jacobpq/soccer-manager/internal/mocks"
)

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		inputBody      models.LoginRequest
		mockBehavior   func(m *mocks.MockAuthService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success - Login OK",
			inputBody: models.LoginRequest{
				Email:    "test@test.com",
				Password: "password123",
			},
			mockBehavior: func(m *mocks.MockAuthService) {
				m.EXPECT().
					Login(gomock.Any(), models.LoginRequest{Email: "test@test.com", Password: "password123"}).
					Return(&models.Session{
						AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.signature",
						RefreshToken: "random-opaque-string",
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "access-123",
		},
		{
			name: "Failure - Invalid Credentials",
			inputBody: models.LoginRequest{
				Email:    "wrong@test.com",
				Password: "wrong_password",
			},
			mockBehavior: func(m *mocks.MockAuthService) {
				m.EXPECT().
					Login(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("invalid credentials"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSvc := mocks.NewMockAuthService(ctrl)
			handler := NewAuthHandler(mockSvc)

			tt.mockBehavior(mockSvc)

			bodyBytes, _ := json.Marshal(tt.inputBody)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(bodyBytes))
			w := httptest.NewRecorder()

			err := handler.Login(w, req)

			if tt.expectedStatus == http.StatusOK {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, w.Code)
				if tt.expectedBody != "" {
					assert.Contains(t, w.Body.String(), tt.expectedBody)
				}
			} else {
				assert.Error(t, err)

				if appErr, ok := err.(*api.AppError); ok {
					assert.Equal(t, tt.expectedStatus, appErr.Status)
				} else {
					t.Errorf("Expected api.AppError, got %T", err)
				}
			}
		})
	}
}
