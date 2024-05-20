package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/db"
	"sync"
)

func main() {
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	connStr := "postgresql://localhost:5432/chirps?sslmode=disable"

	db, err := db.NewDB(connStr)

	if err != nil {
		panic(err)
	}

	defer db.DataBase.Close()

	apiConfig := apiConfig{
		fileserverHits: 0,
		mu:             sync.Mutex{},
		db:             *db,
	}

	mux.Handle("/app/*", http.StripPrefix("/app", apiConfig.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("GET /api/metrics", apiConfig.metricsHandler)
	mux.HandleFunc("GET /admin/metrics", apiConfig.handleAdminMetrics)
	mux.HandleFunc("/api/reset", apiConfig.resetMetrics)

	mux.HandleFunc("POST /api/chirps", apiConfig.chirpsHandler)
	mux.HandleFunc("GET /api/chirps", apiConfig.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiConfig.getChirpByIDHandler)

	mux.HandleFunc("POST /api/users", apiConfig.CreateUserHandler)

	fmt.Println("Server running on port 8080")

	err = server.ListenAndServe()

	if err != nil {
		panic(err)
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, msg)))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Write(data)
	return nil
}

// middlewareMetricsInc increments the fileserverHits counter for each request
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.mu.Lock()
		cfg.fileserverHits += 1
		cfg.mu.Unlock()

		fmt.Println("Middleware Metrics Inc")
		fmt.Println("Fileserver Hits:", cfg.fileserverHits)

		next.ServeHTTP(w, r) // 继续处理请求
	})
}

// healthzHandler returns a simple "OK" response for health checks
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
