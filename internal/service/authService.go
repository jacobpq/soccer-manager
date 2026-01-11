package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/jacobpq/soccer-manager/internal/api"
	"github.com/jacobpq/soccer-manager/internal/config"
	"github.com/jacobpq/soccer-manager/internal/domain/models"
	"github.com/jacobpq/soccer-manager/internal/locales"
	"github.com/jacobpq/soccer-manager/internal/repository"
)

type AuthService interface {
	Register(ctx context.Context, req models.RegisterRequest) error
	Login(ctx context.Context, req models.LoginRequest) (*models.Session, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
}

type authService struct {
	db          *pgxpool.Pool
	userRepo    *repository.UserRepository
	teamRepo    *repository.TeamRepository
	playerRepo  *repository.PlayerRepository
	sessionRepo *repository.SessionRepository
	jwtSecret   []byte
}

func NewAuthService(db *pgxpool.Pool, u *repository.UserRepository, t *repository.TeamRepository, p *repository.PlayerRepository, s *repository.SessionRepository, cfg *config.Config) AuthService {
	return &authService{
		db:          db,
		userRepo:    u,
		teamRepo:    t,
		playerRepo:  p,
		sessionRepo: s,
		jwtSecret:   []byte(cfg.JWTSecret),
	}
}

func (s *authService) Register(ctx context.Context, req models.RegisterRequest) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	user := &models.User{Email: req.Email, Password: string(hashed)}
	if err := s.userRepo.Create(ctx, tx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	team := &models.Team{
		UserID:  user.ID,
		Name:    req.TeamName,
		Country: req.Country,
		Budget:  5000000,
	}
	if err := s.teamRepo.Create(ctx, tx, team); err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	players := s.generateInitialSquad(team.ID)
	if err := s.playerRepo.CreateBatch(ctx, tx, players); err != nil {
		return fmt.Errorf("failed to generate players: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *authService) generateJWT(userID int, tokenType string, duration time.Duration) (string, error) {
	claims := models.AuthClaims{
		UserID:    userID,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *authService) Login(ctx context.Context, req models.LoginRequest) (*models.Session, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, api.ErrUnauthorized(locales.T(ctx, "invalid_credentials"))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, api.ErrUnauthorized(locales.T(ctx, "invalid_credentials"))
	}

	accessToken, _ := s.generateJWT(user.ID, "access", 30*time.Minute)

	refreshToken, _ := s.generateJWT(user.ID, "refresh", 7*24*time.Hour)

	session := &models.Session{
		UserID:           user.ID,
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresAt:  time.Now().Add(15 * time.Minute),
		RefreshExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	session, err := s.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", api.ErrUnauthorized(locales.T(ctx, "invalid_token"))
	}

	newAccessToken, err := s.generateJWT(session.UserID, "access", 30*time.Minute)

	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return newAccessToken, nil
}

func (s *authService) generateInitialSquad(teamID int) []*models.Player {
	var players []*models.Player

	positions := []string{"GK", "GK", "GK"}
	for i := 0; i < 6; i++ {
		positions = append(positions, "DF")
	}
	for i := 0; i < 6; i++ {
		positions = append(positions, "MF")
	}
	for i := 0; i < 5; i++ {
		positions = append(positions, "AT")
	}

	for _, pos := range positions {
		randFirstName := repository.FirstNames[rand.Intn(len(repository.FirstNames))]
		randLastName := repository.LastNames[rand.Intn(len(repository.LastNames))]
		randCountry := repository.Countries[rand.Intn(len(repository.LastNames))]

		players = append(players, &models.Player{
			TeamID:    teamID,
			FirstName: randFirstName,
			LastName:  randLastName,
			Country:   randCountry,
			Age:       rand.Intn(23) + 18,
			Position:  pos,
			Value:     1000000,
		})
	}
	return players
}
