package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jcourtney5/chirpy-server/internal/auth"
)

// DELETE /api/chirps/{chirpID} endpoint handler
func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	// get the chirpID Url param and convert to UUID
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid chirpID", err)
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

	// Get the chirp from the DB
	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		responseWithError(w, http.StatusNotFound, "Invalid chirp", err)
		return
	}

	// make sure user owns the chirp
	if chirp.UserID != userID {
		responseWithError(w, http.StatusForbidden, "Not allowed to delete this chirp", nil)
		return
	}

	// delete the chirp
	err = cfg.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error deleting chirp", err)
		return
	}

	// Send the response
	w.WriteHeader(http.StatusNoContent)
}
