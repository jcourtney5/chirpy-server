package main

import "net/http"

// health endpoint handler function
func handlerHealth(w http.ResponseWriter, r *http.Request) {
	// Supports any HTTP method

	// Add headers
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Set result to 200 OK
	w.WriteHeader(http.StatusOK)

	// Write OK as the response
	_, err := w.Write([]byte(http.StatusText(http.StatusOK)))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
