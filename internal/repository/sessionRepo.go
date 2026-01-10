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
		INSERT INTO sessions (user_id, access_token, refresh_token, access_expires_at, refresh_expires_at) 
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query,
		session.UserID, session.AccessToken, session.RefreshToken,
		session.AccessExpiresAt, session.RefreshExpiresAt)
	return err
}

func (r *SessionRepository) GetUserIDByAccessToken(ctx context.Context, token string) (int, error) {
	var userID int

	query := `SELECT user_id FROM sessions WHERE access_token = $1 AND access_expires_at > $2`
	err := r.db.QueryRow(ctx, query, token, time.Now()).Scan(&userID)
	return userID, err
}

func (r *SessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	var s models.Session

	query := `
		SELECT id, user_id, access_token, refresh_token, access_expires_at, refresh_expires_at 
		FROM sessions WHERE refresh_token = $1 AND refresh_expires_at > $2`

	err := r.db.QueryRow(ctx, query, refreshToken, time.Now()).Scan(
		&s.ID, &s.UserID, &s.AccessToken, &s.RefreshToken, &s.AccessExpiresAt, &s.RefreshExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SessionRepository) RotateAccessToken(ctx context.Context, sessionID int, newAccessToken string, newExpiry time.Time) error {
	query := `UPDATE sessions SET access_token = $1, access_expires_at = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, newAccessToken, newExpiry, sessionID)
	return err
}
