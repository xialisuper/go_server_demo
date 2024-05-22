package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"server/db"
	"server/jwt"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from.env file
	godotenv.Load()
	fmt.Println("Starting server...")

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

	jwtExpireSecStr := os.Getenv("JWT_EXPIRE_SECONDS")
	jwtExpireSec, err := strconv.ParseInt(jwtExpireSecStr, 10, 64)
	if err != nil {
		// 在这里处理转换错误
		panic(err)
	}

	userFreshTokenExpireSecStr := os.Getenv("USER_REFRESH_TOKEN_EXPIRE_SECONDS")
	userFreshTokenExpireSec, err := strconv.ParseInt(userFreshTokenExpireSecStr, 10, 64)
	if err != nil {
		// 在这里处理转换错误
		panic(err)
	}

	apiConfig := ApiConfig{
		fileserverHits:          0,
		mu:                      sync.Mutex{},
		db:                      *db,
		JwtSecret:               os.Getenv("JWT_SECRET"),
		JwtExpireSec:            jwtExpireSec,
		UserFreshTokenExpireSec: userFreshTokenExpireSec,
	}

	mux.Handle("/app/*", http.StripPrefix("/app", apiConfig.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("GET /api/metrics", apiConfig.metricsHandler)
	mux.HandleFunc("GET /admin/metrics", apiConfig.handleAdminMetrics)
	mux.HandleFunc("/api/reset", apiConfig.resetMetrics)

	mux.Handle("POST /api/chirps", apiConfig.authenticationMiddleware(http.HandlerFunc(apiConfig.CreateChirpHandler)))
	mux.HandleFunc("GET /api/chirps", apiConfig.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiConfig.getChirpByIDHandler)

	mux.HandleFunc("POST /api/users", apiConfig.CreateUserHandler)
	//  LOGIN POST /api/login
	mux.HandleFunc("POST /api/login", apiConfig.LoginUserHandler)
	// PUT /api/users
	mux.HandleFunc("PUT /api/users", apiConfig.UpdateUserHandler)
	// POST /api/refresh
	mux.HandleFunc("POST /api/refresh", apiConfig.RefreshTokenHandler)
	// POST /api/revoke
	mux.HandleFunc("POST /api/revoke", apiConfig.RevokeTokenHandler)
	// DELETE /api/chirps/{chirpID}
	mux.Handle("DELETE /api/chirps/{chirpID}", apiConfig.authenticationMiddleware(http.HandlerFunc(apiConfig.deleteChirpByIDHandler)))
	// POST /API/POLKA/WEBHOOKS
	mux.HandleFunc("POST /api/polka/webhooks", apiConfig.PolkaWebhookHandler)

	fmt.Println("Server running on port 8080")

	err = server.ListenAndServe()

	if err != nil {
		panic(err)
	}
}

// respondWithError 函数接收一个 http.ResponseWriter 对象、状态码和消息作为参数，
// 设置响应头的 Content-Type 为 application/json; charset=utf-8，
// 设置状态码并返回错误信息的 JSON 格式。
func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, msg)))

}

// respondWithJSON 函数接收一个 http.ResponseWriter 对象、状态码以及一个任意类型的数据作为参数，
// 设置响应头的 Content-Type 为 application/json; charset=utf-8，
// 设置状态码，将数据转换为 JSON 格式并返回。
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
func (cfg *ApiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.mu.Lock()
		cfg.fileserverHits += 1
		cfg.mu.Unlock()

		fmt.Println("Middleware Metrics Inc")
		fmt.Println("Fileserver Hits:", cfg.fileserverHits)

		next.ServeHTTP(w, r) // 继续处理请求
	})
}

type contextKey string

const userIDKey contextKey = "userID"

// auth middleware function to check if the user is authenticated
func (cfg *ApiConfig) authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// get token from header
		token, err := GetTokenFromHeader(r)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		// validate token

		userIDStr, err := jwt.VerifyJwtToken(token, cfg.JwtSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		// convert the user id from string to int
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		// check if the user exists in the database

		_, err = cfg.db.GetUserByID(userID)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "User not found")
			return
		}

		// Otherwise, continue with the request
		fmt.Println("user authenticated, user id:", userID)

		// save user id in the request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, userIDKey, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// healthzHandler returns a simple "OK" response for health checks
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
