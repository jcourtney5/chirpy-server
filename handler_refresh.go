package main

import (
	"net/http"
	"time"

	"github.com/jcourtney5/chirpy-server/internal/auth"
)

// POST /api/refresh endpoint handler
func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	// Get the refresh token from the authorization header
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	// Get the user from the refresh token (won't return user if expired or revoked)
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "invalid refresh token (expired or revoked)", err)
		return
	}

	// create a new refreshed JWT
	jwtToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error generating JWT token", err)
		return
	}

	// send the response
	responseWithJSON(w, http.StatusOK, response{
		Token: jwtToken,
	})
}

// POST /api/revoke endpoint handler
func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	// Get the refresh token from the authorization header
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	// Revoke the refresh token in the DB
	_, err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't revoke session", err)
	}

	// send no content response (ie success)
	w.WriteHeader(http.StatusNoContent)
}
