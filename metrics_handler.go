package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// metricsHandler returns the number of times Chirpy has been visited
func (a *ApiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits: " + strconv.Itoa(a.fileserverHits)))
}

// resetMetrics resets the fileserverHits counter to 0
func (a *ApiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	a.fileserverHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics reset\n"))
}

// handleAdminMetrics returns a simple HTML page with the number of times Chirpy has been visited
func (cfg *ApiConfig) handleAdminMetrics(w http.ResponseWriter, r *http.Request) {
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
