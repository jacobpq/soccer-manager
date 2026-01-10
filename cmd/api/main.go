package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jacobpq/soccer-manager/internal/config"
	"github.com/jacobpq/soccer-manager/internal/middleware"
	"github.com/jacobpq/soccer-manager/internal/repository"
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

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthCheckHandler(dbPool))

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: middleware.Logger(mux),
	}

	log.Println("Server starting on port " + cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
