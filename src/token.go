package main

import (
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 14).Unix(),
	})
	return token.SignedString(jwtSecretKey)
}

func getBearerToken(token string) string {
	bearer := strings.Split(token, "Bearer ")
	if len(bearer) > 1 {
		return bearer[1]
	}
	return ""
}
