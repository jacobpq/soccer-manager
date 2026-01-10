package main

import (
	"log"
	"net/http"

	"github.com/jacobpq/soccer-manager/internal/config"
	"github.com/jacobpq/soccer-manager/internal/middleware"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "Soccer Manager API is running"}`))
}

func main() {
	cfg := config.LoadConfig()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthCheckHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: middleware.Logger(mux),
	}

	log.Println("Server starting on port " + cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
