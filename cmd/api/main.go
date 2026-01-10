package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jacobpq/soccer-manager/internal/api"
	"github.com/jacobpq/soccer-manager/internal/config"
	"github.com/jacobpq/soccer-manager/internal/handler"
	"github.com/jacobpq/soccer-manager/internal/locales"
	"github.com/jacobpq/soccer-manager/internal/middleware"
	"github.com/jacobpq/soccer-manager/internal/repository"
	"github.com/jacobpq/soccer-manager/internal/service"
)

func healthCheckHandler(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		dbStatus := "up"
		if err := dbPool.Ping(context.Background()); err != nil {
			dbStatus = "down"
			log.Printf("DB health check failed: %v", err)
		}

		response := map[string]string{
			"api_status":      "Soccer Manager API is running",
			"database_status": dbStatus,
		}

		if dbStatus == "down" {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Failed to encode health check response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func main() {
	cfg := config.LoadConfig()

	dbPool, err := repository.InitDB(cfg.DBDSN)
	if err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}
	defer dbPool.Close()

	if err := locales.Init(); err != nil {
		log.Fatal(err)
		log.Fatalf("Failed to init locales: %v", err)

	}

	userRepo := repository.NewUserRepository(dbPool)
	teamRepo := repository.NewTeamRepository()
	playerRepo := repository.NewPlayerRepository()
	sessionRepo := repository.NewSessionRepository(dbPool)

	//service
	authSvc := service.NewAuthService(dbPool, userRepo, teamRepo, playerRepo, sessionRepo)
	teamSvc := service.NewTeamService(dbPool, teamRepo, playerRepo)
	transferSvc := service.NewTransferService(dbPool, playerRepo, teamRepo)

	//handler
	authHandler := handler.NewAuthHandler(authSvc)
	teamHandler := handler.NewTeamHandler(teamSvc)
	transferHandler := handler.NewTransferHandler(transferSvc)

	//middleware
	authMiddleware := middleware.Auth(sessionRepo)

	mux := http.NewServeMux()

	//util
	mux.HandleFunc("GET /health", healthCheckHandler(dbPool))

	//auth
	mux.HandleFunc("POST /register", api.Make(authHandler.Register))
	mux.HandleFunc("POST /login", api.Make(authHandler.Login))
	mux.HandleFunc("POST /refresh", api.Make(authHandler.Refresh))

	//team
	mux.Handle("GET /team", authMiddleware(api.Make(teamHandler.GetMyTeam)))
	mux.Handle("POST /transfer/list", authMiddleware(api.Make(transferHandler.ListPlayer)))
	mux.Handle("POST /transfer/remove", authMiddleware(api.Make(transferHandler.RemovePlayer)))
	mux.Handle("GET /transfer/market", authMiddleware(api.Make(transferHandler.GetMarket)))
	mux.Handle("POST /transfer/buy", authMiddleware(api.Make(transferHandler.BuyPlayer)))
	mux.Handle("PUT /team", authMiddleware(api.Make(teamHandler.UpdateTeam)))
	mux.Handle("PUT /player", authMiddleware(api.Make(teamHandler.UpdatePlayer)))

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: middleware.Logger(middleware.Locale(mux)),
	}

	log.Println("Server starting on port " + cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
