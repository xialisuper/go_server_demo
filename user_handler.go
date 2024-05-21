package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/db"
	"server/jwt"
	"strconv"
	"strings"
)

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

	// login user in database
	user, err = cfg.db.LoginUser(user.Email, user.Password)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// generate JWT token

	if user.Expires == 0 {
		user.Expires = 60 * 60 * 24 * 30 // 30 days
	}

	token, err := jwt.CreateJwtToken(strconv.Itoa(int(user.ID)), user.Expires, cfg.JwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// return JWT token and user data
	resJson := make(map[string]interface{})
	resJson["token"] = token
	resJson["user"] = user

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

	// get jwt from header
	// header "Authorization: Bearer <token>"
	// 从请求头中获取Bearer token
	authHeader := r.Header.Get("Authorization")
	tokenParts := strings.Split(authHeader, " ")
	var token string

	if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
		token = tokenParts[1]
		fmt.Println("Bearer token:", token)
	} else {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
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
