package main

import (
	"server/db"
	"sync"
)

type apiConfig struct {
	fileserverHits int
	mu             sync.Mutex
	db             db.DB
}
