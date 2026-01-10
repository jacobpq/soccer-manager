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
	svc *service.TeamService
}

func NewTeamHandler(svc *service.TeamService) *TeamHandler {
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
