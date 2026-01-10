package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jacobpq/soccer-manager/internal/domain/models"
)

type PlayerRepository struct{}

func NewPlayerRepository() *PlayerRepository {
	return &PlayerRepository{}
}

func (r *PlayerRepository) CreateBatch(ctx context.Context, tx pgx.Tx, players []*models.Player) error {
	query := `
		INSERT INTO players (team_id, first_name, last_name, country, age, position, value, market_value)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	for _, p := range players {
		_, err := tx.Exec(ctx, query,
			p.TeamID, p.FirstName, p.LastName, p.Country,
			p.Age, p.Position, p.Value, 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PlayerRepository) GetByTeamID(ctx context.Context, db *pgxpool.Pool, teamID int) ([]*models.Player, error) {
	query := `
		SELECT id, team_id, first_name, last_name, country, age, position, value, market_value, on_transfer_list 
		FROM players WHERE team_id = $1`

	rows, err := db.Query(ctx, query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make([]*models.Player, 0)

	for rows.Next() {
		var p models.Player
		var marketValue *float64

		err := rows.Scan(
			&p.ID, &p.TeamID, &p.FirstName, &p.LastName, &p.Country,
			&p.Age, &p.Position, &p.Value, &marketValue, &p.OnTransferList,
		)
		if err != nil {
			return nil, err
		}
		if marketValue != nil {
			p.MarketPrice = *marketValue
		}
		players = append(players, &p)
	}
	return players, nil
}

func (r *PlayerRepository) UpdateMarketStatus(ctx context.Context, db *pgxpool.Pool, playerID int, price float64, onList bool) error {
	query := `UPDATE players SET market_value = $1, on_transfer_list = $2 WHERE id = $3`
	_, err := db.Exec(ctx, query, price, onList, playerID)
	return err
}

func (r *PlayerRepository) GetMarketPlayers(ctx context.Context, db *pgxpool.Pool) ([]*models.Player, error) {
	query := `
		SELECT id, team_id, first_name, last_name, country, age, position, value, market_value 
		FROM players WHERE on_transfer_list = true`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make([]*models.Player, 0)

	for rows.Next() {
		var p models.Player
		err := rows.Scan(&p.ID, &p.TeamID, &p.FirstName, &p.LastName, &p.Country, &p.Age, &p.Position, &p.Value, &p.MarketPrice)
		if err != nil {
			return nil, err
		}
		p.OnTransferList = true
		players = append(players, &p)
	}
	return players, nil
}

func (r *PlayerRepository) GetByID(ctx context.Context, db *pgxpool.Pool, playerID int) (*models.Player, error) {
	var p models.Player
	query := `SELECT id, team_id, value, market_value, on_transfer_list FROM players WHERE id = $1`
	err := db.QueryRow(ctx, query, playerID).Scan(&p.ID, &p.TeamID, &p.Value, &p.MarketPrice, &p.OnTransferList)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PlayerRepository) TransferOwnership(ctx context.Context, tx pgx.Tx, playerID int, newTeamID int, newValue float64) error {
	query := `
        UPDATE players 
        SET team_id = $1, value = $2, on_transfer_list = false, market_value = 0 
        WHERE id = $3`
	_, err := tx.Exec(ctx, query, newTeamID, newValue, playerID)
	return err
}

func (r *PlayerRepository) UpdateDetails(ctx context.Context, db *pgxpool.Pool, playerID int, first, last, country string) error {
	query := `UPDATE players SET first_name = $1, last_name = $2, country = $3 WHERE id = $4`
	_, err := db.Exec(ctx, query, first, last, country, playerID)
	return err
}
