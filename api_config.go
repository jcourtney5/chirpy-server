package main

import (
	"sync/atomic"

	"github.com/jcourtney5/chirpy-server/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}
