package main

import (
	"encoding/json"
	"net/http"

	"github.com/jcourtney5/chirpy-server/internal/auth"
	"github.com/jcourtney5/chirpy-server/internal/database"
)

// PUT /api/users endpoint handler
func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	// parse the request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	// Check the authorization header
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	// hash the password
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	// Update the user in the DB
	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Failed to update the user", err)
		return
	}

	// send the response
	responseWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
