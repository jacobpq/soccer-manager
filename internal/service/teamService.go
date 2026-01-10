package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jacobpq/soccer-manager/internal/domain/models"
	"github.com/jacobpq/soccer-manager/internal/repository"
)

type TeamService struct {
	db         *pgxpool.Pool
	teamRepo   *repository.TeamRepository
	playerRepo *repository.PlayerRepository
}

func NewTeamService(db *pgxpool.Pool, t *repository.TeamRepository, p *repository.PlayerRepository) *TeamService {
	return &TeamService{db: db, teamRepo: t, playerRepo: p}
}

type TeamResponse struct {
	Team    *models.Team     `json:"team"`
	Players []*models.Player `json:"players"`
}

func (s *TeamService) GetMyTeam(ctx context.Context, userID int) (*TeamResponse, error) {
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
