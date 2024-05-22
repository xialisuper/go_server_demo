package main

import (
	"encoding/json"
	"log"
	"net/http"
	"server/db"
	"server/jwt"
	"strconv"
	"time"
)

type Event struct {
	Event string `json:"event"`
	Data  Data   `json:"data"`
}

type Data struct {
	UserID int `json:"user_id"`
}

// PolkaWebhookHandler
func (cfg *ApiConfig) PolkaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// get event data from request body
	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	// handle event
	switch event.Event {
	// case "user_created":
	// do something when user created
	case "user.upgraded":

		// update is_chirpy_red to true in  database
		err = cfg.db.UpdateIsChirpyRed(event.Data.UserID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}

	default:
		// do something when event not found
		respondWithError(w, http.StatusNotFound, "event not found")
		return
	}

	// return 204
	respondWithJSON(w, http.StatusNoContent, nil)
}

// RevokeTokenHandler
func (cfg *ApiConfig) RevokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	// 从请求头中获取refresh token
	refreshToken, err := GetTokenFromHeader(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// revoke refresh token in database
	err = cfg.db.RevokeToken(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// return 204
	respondWithJSON(w, http.StatusNoContent, nil)
}

// RefreshTokenHandler
func (cfg *ApiConfig) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	// 从请求头中获取refresh token
	refreshToken, err := GetTokenFromHeader(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// check refresh token in database
	userID, err := cfg.db.CheckRefreshTokenIsValid(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// generate JWT token
	token, err := jwt.CreateJwtToken(strconv.Itoa(userID), cfg.JwtSecret, cfg.JwtExpireSec)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// return JWT token
	resJson := make(map[string]string)
	resJson["token"] = token
	respondWithJSON(w, http.StatusOK, resJson)

}

func (cfg *ApiConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	// get user data from request body
	var user db.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// create user in database
	user, err = cfg.db.CreateUser(user.Email, user.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, user)

}

func (cfg *ApiConfig) LoginUserHandler(w http.ResponseWriter, r *http.Request) {

	// get user data from request body
	var user db.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// check user password in database
	user, err = cfg.db.LoginUser(user.Email, user.Password)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// generate JWT token
	token, err := jwt.CreateJwtToken(strconv.Itoa(int(user.ID)), cfg.JwtSecret, cfg.JwtExpireSec)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// generate refresh token
	refreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// return JWT token and user data
	resJson := make(map[string]interface{})
	resJson["token"] = token
	resJson["user"] = user
	resJson["refresh_token"] = refreshToken

	// save refresh token in database
	// refresh token expiration time to set to 60 days by default

	expire_time := time.Now().Add(time.Duration(cfg.UserFreshTokenExpireSec) * time.Second)
	err = cfg.db.SaveToken(user.ID, refreshToken, expire_time)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, resJson)

}

// UpdateUserHandler
func (cfg *ApiConfig) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {

	// get user data from request body
	var user db.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		log.Println(err)
		return
	}

	// 从请求头中获取Bearer token
	token, err := GetTokenFromHeader(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// parse JWT token  claims is user id
	claims, err := jwt.VerifyJwtToken(token, cfg.JwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// convert user id from string to int
	userID, err := strconv.Atoi(claims)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	// update user in database
	user, err = cfg.db.UpdateUser(userID, user.Email, user.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)

}
