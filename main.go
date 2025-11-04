package main

import (
	"fmt"
	"net/http"
)

// to run: go build -o out && ./out

func main() {
	// create the server
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Starting HTTP server at port 8080")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
