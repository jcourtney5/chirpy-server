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
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

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
