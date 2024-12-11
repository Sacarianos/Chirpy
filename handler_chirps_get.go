package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("Error retrieving chirps: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve chirps", err)
		return
	}

	// Map database chirps to API response format
	var response []Chirp
	for _, chirp := range chirps {
		response = append(response, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	// Send the response
	respondWithJSON(w, http.StatusOK, response)

}

func (cfg *apiConfig) handleGetChirpById(w http.ResponseWriter, r *http.Request) {

	// Extract the chirpID from the URL path
	chirpID := strings.TrimPrefix(r.URL.Path, "/api/chirps/")

	// Parse chirpID as UUID
	id, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID format", err)
		return
	}
	// Query the chirp from the database
	chirp, err := cfg.db.GetChirpById(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Chirp not found", err)
			return
		}
		log.Printf("Error retrieving chirp: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve chirp", err)
		return
	}

	// Respond with the chirp as JSON
	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		UserID:    chirp.UserID,
		Body:      chirp.Body,
	})

}
