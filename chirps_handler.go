package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"server/db"
	"strconv"
	"strings"
)

func (cfg *ApiConfig) getChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
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

func (cfg *ApiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	// Get all chirps from the database
	chirps, err := cfg.db.GetChirps()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 200 OK
	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *ApiConfig) CreateChirpHandler(w http.ResponseWriter, r *http.Request) {

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

	//  use r.context.Value("userID") instead of parsing the JWT token again
	userID := r.Context().Value(userIDKey).(int)


	// Save the chirp to the database
	newChirp, err := cfg.db.CreateChirp(validatedChirp.Body, userID)

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
