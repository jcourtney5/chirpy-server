package main

import (
	"fmt"
	"net/http"
)

// to run: go build -o out && ./out

func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	// add our handlers
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot))))
	mux.HandleFunc("/healthz", handlerHealth)

	// create our server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("Starting HTTP server at port 8080")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}

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
