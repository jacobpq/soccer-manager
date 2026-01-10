package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func Make(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			if e, ok := err.(*AppError); ok {
				if e.Err != nil {
					log.Printf("Internal Error: %v", e.Err)
				}

				WriteError(w, e.Status, e.Msg)
				return
			}

			log.Printf("Unhandled Error: %v", err)
			WriteError(w, http.StatusInternalServerError, "Internal Server Error")
		}
	}
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
