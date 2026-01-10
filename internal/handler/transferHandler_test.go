package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/jacobpq/soccer-manager/internal/api"
	"github.com/jacobpq/soccer-manager/internal/middleware"
	"github.com/jacobpq/soccer-manager/internal/mocks"
)

func TestTransferHandler_BuyPlayer(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		inputBody      map[string]interface{}
		mockBehavior   func(m *mocks.MockTransferService)
		expectedStatus int
	}{
		{
			name:   "Success - Player Bought",
			userID: 55,
			inputBody: map[string]interface{}{
				"player_id": 100,
			},
			mockBehavior: func(m *mocks.MockTransferService) {
				m.EXPECT().
					BuyPlayer(gomock.Any(), 55, 100).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Failure - Insufficient Funds",
			userID: 55,
			inputBody: map[string]interface{}{
				"player_id": 100,
			},
			mockBehavior: func(m *mocks.MockTransferService) {
				m.EXPECT().
					BuyPlayer(gomock.Any(), 55, 100).
					Return(errors.New("insufficient funds"))
			},

			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Failure - Invalid JSON",
			userID: 55,
			inputBody: map[string]interface{}{
				"player_id": "this-should-be-int",
			},
			mockBehavior: func(m *mocks.MockTransferService) {
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSvc := mocks.NewMockTransferService(ctrl)
			handler := NewTransferHandler(mockSvc)

			tt.mockBehavior(mockSvc)

			bodyBytes, _ := json.Marshal(tt.inputBody)
			req := httptest.NewRequest(http.MethodPost, "/transfer/buy", bytes.NewBuffer(bodyBytes))
			w := httptest.NewRecorder()

			ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.userID)
			req = req.WithContext(ctx)

			err := handler.BuyPlayer(w, req)

			if tt.expectedStatus == http.StatusOK {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, w.Code)
			} else {
				assert.Error(t, err)
				if appErr, ok := err.(*api.AppError); ok {
					assert.Equal(t, tt.expectedStatus, appErr.Status)
				}
			}
		})
	}
}
