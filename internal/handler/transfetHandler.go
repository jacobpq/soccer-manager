package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jacobpq/soccer-manager/internal/api"
	"github.com/jacobpq/soccer-manager/internal/locales"
	"github.com/jacobpq/soccer-manager/internal/middleware"
	"github.com/jacobpq/soccer-manager/internal/service"
)

type TransferHandler struct {
	svc *service.TransferService
}

func NewTransferHandler(svc *service.TransferService) *TransferHandler {
	return &TransferHandler{svc: svc}
}

type ListRequest struct {
	PlayerID int     `json:"player_id"`
	Price    float64 `json:"price"`
}

type BuyRequest struct {
	PlayerID int `json:"player_id"`
}

func (h *TransferHandler) ListPlayer(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := ctx.Value(middleware.UserIDKey).(int)

	var req ListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return api.ErrBadRequest(locales.T(ctx, "invalid_json"))
	}

	if err := h.svc.ListPlayer(ctx, userID, req.PlayerID, req.Price); err != nil {
		return api.ErrBadRequest(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(map[string]string{
		"status": locales.T(ctx, "player_listed"),
	})
}

func (h *TransferHandler) RemovePlayer(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := ctx.Value(middleware.UserIDKey).(int)

	var req BuyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return api.ErrBadRequest(locales.T(ctx, "invalid_json"))
	}

	if err := h.svc.RemoveFromList(ctx, userID, req.PlayerID); err != nil {
		return api.ErrBadRequest(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(map[string]string{
		"status": locales.T(ctx, "player_removed_from_list"),
	})
}

func (h *TransferHandler) GetMarket(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	players, err := h.svc.GetMarket(ctx)
	if err != nil {
		return api.NewError(err, http.StatusInternalServerError, locales.T(ctx, "market_fetch_fail"))
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(players)
}

func (h *TransferHandler) BuyPlayer(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := ctx.Value(middleware.UserIDKey).(int)

	var req BuyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return api.ErrBadRequest(locales.T(ctx, "invalid_json"))
	}

	if err := h.svc.BuyPlayer(ctx, userID, req.PlayerID); err != nil {
		msg := err.Error()
		return api.ErrBadRequest(msg)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(map[string]string{
		"status": locales.T(ctx, "transfer_success"),
	})
}
