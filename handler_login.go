package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jcourtney5/chirpy-server/internal/auth"
)

const defaultExpiresInSeconds = 3600 // one hour

// POST /api/login endpoint handler
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
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

	// If ExpiresInSeconds is missing or over the default, use the default
	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds > defaultExpiresInSeconds {
		params.ExpiresInSeconds = defaultExpiresInSeconds
	}
	expiresIn := time.Duration(params.ExpiresInSeconds) * time.Second

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
	jwtToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiresIn)

	// send the response
	responseWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: jwtToken,
	})
}
