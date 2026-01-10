package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jacobpq/soccer-manager/internal/api"
	"github.com/jacobpq/soccer-manager/internal/domain/models"
	"github.com/jacobpq/soccer-manager/internal/locales"
	"github.com/jacobpq/soccer-manager/internal/repository"
)

type TeamService interface {
	GetMyTeam(ctx context.Context, userID int) (*TeamResponse, error)
	UpdateTeam(ctx context.Context, userID int, name, country *string) error
	UpdatePlayer(ctx context.Context, userID, playerID int, first, last, country *string) error
}

type teamService struct {
	db         *pgxpool.Pool
	teamRepo   *repository.TeamRepository
	playerRepo *repository.PlayerRepository
}

func NewTeamService(db *pgxpool.Pool, t *repository.TeamRepository, p *repository.PlayerRepository) TeamService {
	return &teamService{db: db, teamRepo: t, playerRepo: p}
}

type TeamResponse struct {
	Team    *models.Team     `json:"team"`
	Players []*models.Player `json:"players"`
}

func (s *teamService) GetMyTeam(ctx context.Context, userID int) (*TeamResponse, error) {
	team, err := s.teamRepo.GetByUserID(ctx, s.db, userID)
	if err != nil {
		return nil, err
	}

	players, err := s.playerRepo.GetByTeamID(ctx, s.db, team.ID)
	if err != nil {
		return nil, err
	}

	var totalValue float64
	for _, p := range players {
		totalValue += p.Value
	}
	team.Value = totalValue

	return &TeamResponse{
		Team:    team,
		Players: players,
	}, nil
}

func (s *teamService) UpdateTeam(ctx context.Context, userID int, name, country *string) error {
	team, err := s.teamRepo.GetByUserID(ctx, s.db, userID)
	if err != nil {
		return err
	}

	if name != nil {
		team.Name = *name
	}
	if country != nil {
		team.Country = *country
	}

	return s.teamRepo.UpdateDetails(ctx, s.db, team.ID, team.Name, team.Country)
}

func (s *teamService) UpdatePlayer(ctx context.Context, userID, playerID int, first, last, country *string) error {
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

	if first != nil {
		player.FirstName = *first
	}
	if last != nil {
		player.LastName = *last
	}
	if country != nil {
		player.Country = *country
	}

	return s.playerRepo.UpdateDetails(ctx, s.db, playerID, player.FirstName, player.LastName, player.Country)
}
