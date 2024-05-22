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



