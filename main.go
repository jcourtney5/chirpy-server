package main

import (
	"database/sql"
	"fmt"
	"github.com/jcourtney5/chirpy-server/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

// sql: postgres://jessecourtney:@localhost:5432/chirpy?sslmode=disable
// goose postgres "postgres://jessecourtney:@localhost:5432/chirpy?sslmode=disable" up
// goose postgres "postgres://jessecourtney:@localhost:5432/chirpy?sslmode=disable" down

// to run: go build -o out && ./out
//     or: go run .

func main() {
	const filepathRoot = "."
	const port = "8080"

	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// load our env variables
	dbURL := os.Getenv("DB_URL")

	// connect to the db
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
	}

	mux := http.NewServeMux()

	// main file server handler
	fileServerHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServerHandler))

	// api handlers
	mux.HandleFunc("GET /api/healthz", handlerHealth)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	// admin handlers
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	// create our server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("Starting HTTP server at port 8080")
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
