package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jacobpq/soccer-manager/internal/api"
	"github.com/jacobpq/soccer-manager/internal/domain/models"
	"github.com/jacobpq/soccer-manager/internal/locales"
	"github.com/jacobpq/soccer-manager/internal/repository"
	"github.com/jacobpq/soccer-manager/internal/service"
)

type AuthHandler struct {
	svc *service.AuthService
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return api.ErrBadRequest(locales.T(ctx, "invalid_json"))
	}

	if err := req.Validate(); err != nil {
		return api.ErrBadRequest(locales.T(ctx, err.Error()))
	}

	if err := h.svc.Register(r.Context(), req); err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			return api.NewError(nil, http.StatusConflict, locales.T(ctx, "email_exists"))
		}

		return api.ErrInternal(err)
	}

	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(map[string]string{"message": locales.T(ctx, "user_created")})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return api.ErrBadRequest("Invalid JSON body")
	}

	if err := req.Validate(); err != nil {
		return api.ErrBadRequest(locales.T(ctx, err.Error()))
	}

	session, err := h.svc.Login(r.Context(), req)
	if err != nil {
		return api.ErrUnauthorized(locales.T(r.Context(), "invalid_credentials"))
	}

	return json.NewEncoder(w).Encode(map[string]string{
		"access_token":  session.AccessToken,
		"refresh_token": session.RefreshToken,
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) error {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return api.ErrBadRequest("Invalid JSON")
	}

	newAccessToken, err := h.svc.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		return api.ErrUnauthorized(locales.T(r.Context(), "invalid_refresh_token"))
	}

	return json.NewEncoder(w).Encode(map[string]string{
		"access_token": newAccessToken,
	})
}
