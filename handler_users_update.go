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

func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode the request body
	var req request
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// insert user into the DB
	user, err := cfg.db.UpdateUser(context.Background(), database.UpdateUserParams{
		Email:          req.Email,
		HashedPassword: hashedPassword,
		ID:             userId,
	})
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	type UserResponse struct {
		ID        uuid.UUID `json:"id"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	resp := UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	respondWithJSON(w, http.StatusOK, resp)

}
