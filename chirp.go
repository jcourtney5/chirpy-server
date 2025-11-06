package main

import (
	"encoding/json"
	"net/http"
)

// /api/validate_chip endpoint handler
func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type validResponse struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		responseWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	responseWithJSON(w, http.StatusOK, validResponse{Valid: true})
}
