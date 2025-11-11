package main

import (
	"encoding/json"
	"net/http"

	"github.com/jcourtney5/chirpy-server/internal/auth"
)

// POST /api/login endpoint handler
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	// send the response
	responseWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
