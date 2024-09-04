package token

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTMaker interface {
	CreateToken(string, string) (string, error)
	Parse(string) (*jwt.Token, error)
}

type jwtMaker struct {
	secret []byte
}

func GetBearerToken(token string) string {
	bearer := strings.Split(token, "Bearer ")
	if len(bearer) > 1 {
		return bearer[1]
	}
	return ""
}

func NewJWTMaker(secret []byte) *jwtMaker {
	return &jwtMaker{secret: secret}
}

func (j *jwtMaker) CreateToken(username string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24 * 14).Unix(),
	})
	return token.SignedString(j.secret)
}

func (j *jwtMaker) Parse(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})
}
