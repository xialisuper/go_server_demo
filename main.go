package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	apiConfig := apiConfig{
		fileserverHits: 0,
	}

	mux.Handle("/app/*", http.StripPrefix("/app", apiConfig.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.Handle("/healthz", http.HandlerFunc(healthzHandler))
	mux.Handle("/metrics", http.HandlerFunc(apiConfig.metricsHandler))
	mux.Handle("/reset", http.HandlerFunc(apiConfig.resetMetrics))

	err := server.ListenAndServe()

	if err != nil {
		panic(err)
	}
}
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

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (a *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits: " + strconv.Itoa(a.fileserverHits)))
}

func (a *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	a.fileserverHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics reset\n"))
}
