package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jacobpq/soccer-manager/internal/domain/models"
)

type SessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, session *models.Session) error {
	query := `
		INSERT INTO sessions (user_id, refresh_token, refresh_expires_at) 
		VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query,
		session.UserID, session.RefreshToken, session.RefreshExpiresAt)
	return err
}

func (r *SessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	var s models.Session

	query := `
		SELECT id, user_id, refresh_token, refresh_expires_at 
		FROM sessions WHERE refresh_token = $1 AND refresh_expires_at > $2`

	err := r.db.QueryRow(ctx, query, refreshToken, time.Now()).Scan(
		&s.ID, &s.UserID, &s.RefreshToken, &s.RefreshExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
