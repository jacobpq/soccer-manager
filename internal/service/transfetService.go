package service

import (
	"context"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jacobpq/soccer-manager/internal/api"
	"github.com/jacobpq/soccer-manager/internal/domain/models"
	"github.com/jacobpq/soccer-manager/internal/locales"
	"github.com/jacobpq/soccer-manager/internal/repository"
)

type TransferService struct {
	db         *pgxpool.Pool
	playerRepo *repository.PlayerRepository
	teamRepo   *repository.TeamRepository
}

func NewTransferService(db *pgxpool.Pool, p *repository.PlayerRepository, t *repository.TeamRepository) *TransferService {
	return &TransferService{db: db, playerRepo: p, teamRepo: t}
}

func (s *TransferService) ListPlayer(ctx context.Context, userID, playerID int, price float64) error {
	player, err := s.playerRepo.GetByID(ctx, s.db, playerID)
	if err != nil {
		return api.ErrNotFound(locales.T(ctx, "player_not_found"))
	}

	team, err := s.teamRepo.GetByUserID(ctx, s.db, userID)
	if err != nil {
		return api.ErrNotFound(locales.T(ctx, "team_not_found"))
	}

	if player.TeamID != team.ID {
		return api.ErrNotFound(locales.T(ctx, "do_not_own_player"))
	}

	return s.playerRepo.UpdateMarketStatus(ctx, s.db, playerID, price, true)
}

func (s *TransferService) RemoveFromList(ctx context.Context, userID, playerID int) error {
	player, err := s.playerRepo.GetByID(ctx, s.db, playerID)
	if err != nil {
		return api.ErrNotFound(locales.T(ctx, "player_not_found"))

	}

	team, err := s.teamRepo.GetByUserID(ctx, s.db, userID)
	if err != nil {
		return api.ErrNotFound(locales.T(ctx, "team_not_found"))
	}

	if player.TeamID != team.ID {
		return api.ErrNotFound(locales.T(ctx, "do_not_own_player"))
	}

	if !player.OnTransferList {
		return api.ErrNotFound(locales.T(ctx, "player_not_for_sale"))
	}

	return s.playerRepo.UpdateMarketStatus(ctx, s.db, playerID, 0, false)
}

func (s *TransferService) GetMarket(ctx context.Context) ([]*models.Player, error) {
	return s.playerRepo.GetMarketPlayers(ctx, s.db)
}

func (s *TransferService) BuyPlayer(ctx context.Context, buyerUserID, playerID int) error {
	buyerTeam, err := s.teamRepo.GetByUserID(ctx, s.db, buyerUserID)
	if err != nil {
		return api.ErrNotFound(locales.T(ctx, "buyer_team_not_found"))
	}

	player, err := s.playerRepo.GetByID(ctx, s.db, playerID)
	if err != nil {
		return api.ErrNotFound(locales.T(ctx, "player_not_found"))
	}

	if !player.OnTransferList {
		return api.ErrNotFound(locales.T(ctx, "player_not_for_sale"))
	}

	if buyerTeam.Budget < player.MarketPrice {
		return api.ErrNotFound(locales.T(ctx, "insufficient_funds"))
	}

	if player.TeamID == buyerTeam.ID {
		return api.ErrNotFound(locales.T(ctx, "own_player_buy"))
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.teamRepo.UpdateBudget(ctx, tx, buyerTeam.ID, -player.MarketPrice); err != nil {
		return err
	}
	if err := s.teamRepo.UpdateBudget(ctx, tx, player.TeamID, player.MarketPrice); err != nil {
		return err
	}

	rand.Seed(time.Now().UnixNano())
	factor := 1.1 + rand.Float64()*0.9
	newValue := player.Value * factor

	if err := s.playerRepo.TransferOwnership(ctx, tx, playerID, buyerTeam.ID, newValue); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
