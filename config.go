package main

import (
	"server/db"
	"sync"
)

type ApiConfig struct {
	fileserverHits int
	mu             sync.Mutex
	db             db.DB
	JwtSecret      string
	JwtExpireSec   int64
	UserFreshTokenExpireSec int64
}
