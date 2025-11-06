package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func responseWithError(w http.ResponseWriter, statusCode int, message string, err error) {
	if err != nil {
		log.Println(err)
	}
	if statusCode > 499 {
		log.Printf("Request failed with 5XX status code: %s\n", message)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	responseWithJSON(w, statusCode, errorResponse{
		Error: message,
	})
}

func responseWithJSON(w http.ResponseWriter, statusCode int, jsonData interface{}) {
	data, err := json.Marshal(jsonData)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}
