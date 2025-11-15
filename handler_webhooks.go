package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jcourtney5/chirpy-server/internal/auth"
)

// POST /api/polka/webhooks endpoint handler
func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	// parse the request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	// Check the authorization header for the
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Couldn't find ApiKey", err)
		return
	}

	// validate the apiKey
	if apiKey != cfg.polkaKey {
		responseWithError(w, http.StatusUnauthorized, "Invalid ApiKey", nil)
		return
	}

	// don't handle other events yet
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// handle user.upgraded event, upgrade user to chirpy red
	_, err = cfg.db.UpgradeToChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			responseWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		responseWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	// send no content response (ie success)
	w.WriteHeader(http.StatusNoContent)
}
