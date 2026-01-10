package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jacobpq/soccer-manager/internal/api"
	"github.com/jacobpq/soccer-manager/internal/locales"
	"github.com/jacobpq/soccer-manager/internal/middleware"
	"github.com/jacobpq/soccer-manager/internal/service"
)

type TeamHandler struct {
	svc service.TeamService
}

type UpdateTeamRequest struct {
	Name    *string `json:"name"`
	Country *string `json:"country"`
}

type UpdatePlayerRequest struct {
	PlayerID  int     `json:"player_id"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Country   *string `json:"country"`
}

func NewTeamHandler(svc service.TeamService) *TeamHandler {
	return &TeamHandler{svc: svc}
}

func (h *TeamHandler) GetMyTeam(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		return api.ErrUnauthorized(locales.T(ctx, "unauthorized"))
	}

	resp, err := h.svc.GetMyTeam(r.Context(), userID)
	if err != nil {
		return api.NewError(err, http.StatusNotFound, locales.T(ctx, "team_not_found"))
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resp)
}

func (h *TeamHandler) UpdateTeam(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := ctx.Value(middleware.UserIDKey).(int)

	var req UpdateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return api.ErrBadRequest(locales.T(ctx, "invalid_json"))
	}

	if err := h.svc.UpdateTeam(ctx, userID, req.Name, req.Country); err != nil {
		return api.ErrInternal(err)
	}

	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(map[string]string{
		"status": locales.T(ctx, "team_updated"),
	})
}

func (h *TeamHandler) UpdatePlayer(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := ctx.Value(middleware.UserIDKey).(int)

	var req UpdatePlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return api.ErrBadRequest(locales.T(ctx, "invalid_json"))
	}

	if err := h.svc.UpdatePlayer(ctx, userID, req.PlayerID, req.FirstName, req.LastName, req.Country); err != nil {
		return api.ErrBadRequest(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(map[string]string{
		"status": locales.T(ctx, "player_updated"),
	})
}
