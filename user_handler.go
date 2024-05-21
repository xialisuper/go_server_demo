package main

import (
	"encoding/json"
	"net/http"
	"server/db"
)

func (cfg *apiConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {

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

func (cfg *apiConfig) LoginUserHandler(w http.ResponseWriter, r *http.Request) {

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

	respondWithJSON(w, http.StatusOK, user)

}
