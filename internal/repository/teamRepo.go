package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jacobpq/soccer-manager/internal/domain/models"
)

type TeamRepository struct{}

func NewTeamRepository() *TeamRepository {
	return &TeamRepository{}
}

func (r *TeamRepository) Create(ctx context.Context, tx pgx.Tx, team *models.Team) error {
	query := `
		INSERT INTO teams (user_id, name, country, budget) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id`

	err := tx.QueryRow(ctx, query, team.UserID, team.Name, team.Country, team.Budget).Scan(&team.ID)
	return err
}

func (r *TeamRepository) GetByUserID(ctx context.Context, db *pgxpool.Pool, userID int) (*models.Team, error) {
	var team models.Team
	query := `SELECT id, user_id, name, country, budget FROM teams WHERE user_id = $1`
	err := db.QueryRow(ctx, query, userID).Scan(&team.ID, &team.UserID, &team.Name, &team.Country, &team.Budget)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *TeamRepository) UpdateBudget(ctx context.Context, tx pgx.Tx, teamID int, amount float64) error {
	query := `UPDATE teams SET budget = budget + $1 WHERE id = $2`
	_, err := tx.Exec(ctx, query, amount, teamID)
	return err
}
