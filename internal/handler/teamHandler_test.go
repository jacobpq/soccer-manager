package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/jacobpq/soccer-manager/internal/api"
	"github.com/jacobpq/soccer-manager/internal/domain/models"
	"github.com/jacobpq/soccer-manager/internal/middleware"
	"github.com/jacobpq/soccer-manager/internal/mocks"
	"github.com/jacobpq/soccer-manager/internal/service"
)

func TestTeamHandler_GetMyTeam(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		mockBehavior   func(m *mocks.MockTeamService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Success - Found Team",
			userID: 1,
			mockBehavior: func(m *mocks.MockTeamService) {
				m.EXPECT().
					GetMyTeam(gomock.Any(), 1).Return(&service.TeamResponse{
					Team: &models.Team{Name: "Real Madrid"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Real Madrid",
		},
		{
			name:   "Failure - Team Not Found",
			userID: 1,
			mockBehavior: func(m *mocks.MockTeamService) {
				m.EXPECT().
					GetMyTeam(gomock.Any(), 1).
					Return(nil, errors.New("team not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSvc := mocks.NewMockTeamService(ctrl)
			handler := NewTeamHandler(mockSvc)

			tt.mockBehavior(mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/team", nil)
			w := httptest.NewRecorder()

			ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.userID)
			req = req.WithContext(ctx)

			err := handler.GetMyTeam(w, req)

			if tt.expectedStatus == http.StatusOK {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			} else {
				assert.Error(t, err)
				if appErr, ok := err.(*api.AppError); ok {
					assert.Equal(t, tt.expectedStatus, appErr.Status)
				}
			}
		})
	}
}
