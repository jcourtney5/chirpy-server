package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jcourtney5/chirpy-server/internal/auth"
	"github.com/jcourtney5/chirpy-server/internal/database"
)

// POST /api/login endpoint handler
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	// parse the request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	// find the user
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	// verify the hashed password
	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		responseWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	// create the JWT
	jwtToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error generating JWT token", err)
	}

	// create a refresh token and save to DB with 60 day expiration
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error generating refresh token", err)
		return
	}
	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(24 * 60 * time.Hour),
	})
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error saving refresh token", err)
		return
	}

	// send the response
	responseWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        jwtToken,
		RefreshToken: refreshToken,
	})
}
