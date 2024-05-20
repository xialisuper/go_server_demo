package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"server/db"
	"strconv"
	"strings"
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

	fmt.Println("Server running on port 8080")

	err = server.ListenAndServe()

	if err != nil {
		panic(err)
	}
}

func (cfg *apiConfig) getChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Get the chirp ID from the URL /api/chirps/{chirpID}
	chirpID := r.PathValue("chirpID")

	// Get the chirp from the database
	chirpIDInt, err := strconv.Atoi(chirpID)
	if err != nil {
		// 错误处理
		respondWithError(w, http.StatusNotFound, "invalid chirp ID")
		return
	}

	chirp, err := cfg.db.GetChirpByID(chirpIDInt)

	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	// 200 OK
	respondWithJSON(w, http.StatusOK, chirp)
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	// Get all chirps from the database
	chirps, err := cfg.db.GetChirps()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 200 OK
	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) chirpsHandler(w http.ResponseWriter, r *http.Request) {

	var chirp db.Chirp
	err := json.NewDecoder(r.Body).Decode(&chirp)

	if err != nil {
		// 500 Internal Server Error
		respondWithError(w, http.StatusInternalServerError, "something went wrong")
		return

	}

	validatedChirp, err := validateChirp(chirp)

	if err != nil {
		// 400 Bad Request
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	// Save the chirp to the database
	newChirp, err := cfg.db.CreateChirp(validatedChirp.Body)

	if err != nil {

		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 200 OK
	respondWithJSON(w, http.StatusOK, newChirp)
}

// validateChirp validates the chirp and returns a cleaned version of the chirp
func validateChirp(chirp db.Chirp) (validateChirp db.Chirp, err error) {

	// Check if chirp is too long
	if len(chirp.Body) > 140 {
		err = errors.New("chirp is too long")
		return db.Chirp{}, err
	}

	// Replace bad words with "****"
	cleanedBody := replaceBadWords(chirp.Body)

	validatedChirp := db.Chirp{Body: cleanedBody}

	return validatedChirp, nil
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

// replaceBadWords replaces bad words with "****"
func replaceBadWords(chirp string) string {
	wordsToReplace := []string{"kerfuffle", "sharbert", "fornax"}
	replacement := "****"

	// Create a regular expression to match words irrespective of case, excluding adjacent punctuation
	re := regexp.MustCompile(`(?i)\b(` + strings.Join(wordsToReplace, "|") + `)\b`)

	// Replace matching words with "****"
	modifiedChirp := re.ReplaceAllString(chirp, replacement)

	return modifiedChirp
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

// handleAdminMetrics returns a simple HTML page with the number of times Chirpy has been visited
func (cfg *apiConfig) handleAdminMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	//return a html template
	visitCount := cfg.fileserverHits
	htmlContent := fmt.Sprintf(`
	<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
	</html>
`, visitCount)
	_, err := w.Write([]byte(htmlContent))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// healthzHandler returns a simple "OK" response for health checks
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// metricsHandler returns the number of times Chirpy has been visited
func (a *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits: " + strconv.Itoa(a.fileserverHits)))
}

// resetMetrics resets the fileserverHits counter to 0
func (a *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	a.fileserverHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics reset\n"))
}
