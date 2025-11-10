package main

import (
	"net/http"
)

// GET /api/chirps endpoint handler
func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	// Get all the chirps in the DB
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Failed to get chirps", err)
		return
	}

	// convert DB results to our Chirp struct
	chirps := make([]Chirp, 0, len(dbChirps))
	for _, chirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	// Send the response
	responseWithJSON(w, http.StatusOK, chirps)
}
