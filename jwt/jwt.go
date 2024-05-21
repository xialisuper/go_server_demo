package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJwtToken(userId string, expireTimeInSec int64, jwtSecret string) (string, error) {

	// expireTime
	expire := time.Now().Add(time.Duration(expireTimeInSec) * time.Second)
	// generate token

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject: userId,
		// expire time in seconds from now
		ExpiresAt: jwt.NewNumericDate(expire),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "go-server",
	},
	)

	// sign token with secret key

	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJwtToken(tokenString string, jwtSecret string) (string, error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)

	if !ok {
		return "", err
	}

	return claims.Subject, nil
}
