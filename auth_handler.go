package main

import (
	"errors"
	"net/http"
	"strings"
)

// getUserFromToken
// header "Authorization: Bearer <token>"
func GetTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	tokenParts := strings.Split(authHeader, " ")
	var token string

	if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
		token = tokenParts[1]
		return token, nil

	}

	return "", errors.New("invalid token")
}

func GetPolkaApiKeyFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	keyParts := strings.Split(authHeader, " ")
	var apiKey string

	if len(keyParts) == 2 && keyParts[0] == "ApiKey" {
		apiKey = keyParts[1]
		return apiKey, nil

	}

	return "", errors.New("invalid api key")
}
