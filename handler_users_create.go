package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Sacarians/chirpy/internal/auth"
	"github.com/Sacarians/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode the request body
	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// insert user into the DB
	user, err := cfg.db.CreateUser(context.Background(), database.CreateUserParams{
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	type UserResponse struct {
		ID          uuid.UUID `json:"id"`
		Email       string    `json:"email"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
	}

	var chirpyRed bool
	if user.IsChirpyRed.Valid {
		chirpyRed = user.IsChirpyRed.Bool
	} else {
		chirpyRed = false // or another default value based on your application logic
	}

	resp := UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		IsChirpyRed: chirpyRed,
	}

	respondWithJSON(w, http.StatusCreated, resp)

}
