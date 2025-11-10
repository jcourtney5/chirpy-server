package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

// GET /api/chirps/{chirpID} endpoint handler
func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	// get the chirpID Url param and convert to UUID
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid hirpID", err)
		return
	}

	// Get the chirp from the DB
	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		// Check if not found vs general error
		if errors.Is(err, sql.ErrNoRows) {
			responseWithError(w, http.StatusNotFound, "Failed to get chirp", err)
		} else {
			responseWithError(w, http.StatusInternalServerError, "Failed to get chirp", err)
		}
		return
	}

	// Send the response
	responseWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	})
}

// GET /api/chirps endpoint handler
func (cfg *apiConfig) handlerChirpsGetAll(w http.ResponseWriter, r *http.Request) {
	// Get all the chirps in the DB
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Failed to get chirps", err)
		return
	}

	// convert DB results to our Chirp struct
	chirps := make([]Chirp, 0, len(dbChirps))
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	// Send the response
	responseWithJSON(w, http.StatusOK, chirps)
}
