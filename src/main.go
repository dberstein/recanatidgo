package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	ginratelimit "github.com/ljahier/gin-ratelimit"

	_ "github.com/mattn/go-sqlite3"
)

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Email    string `json:"email"`
}

var jwtSecretKey = []byte("your-secret-key")
var db *sql.DB
var tb *ginratelimit.TokenBucket

func init() {
	db = GetDb()
	// Create a new token bucket rate limiter
	tb = ginratelimit.NewTokenBucket(5, 1*time.Minute) // 5 requests per minute
}

func GetDb() *sql.DB {
	db, err := sql.Open("sqlite3", "x.db")
	if err != nil {
		panic(err)
	}

	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS users (username TEXT PRIMARY KEY, pwhash TEXT, email TEXT, role TEXT)",
	); err != nil {
		panic(err)
	}

	return db
}

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 14).Unix(),
	})
	return token.SignedString(jwtSecretKey)
}

func GetBearerToken(token string) string {
	bearer := strings.Split(token, "Bearer ")
	if len(bearer) > 1 {
		return bearer[1]
	}
	return ""
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := GetBearerToken(c.GetHeader("Authorization"))
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecretKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("username", claims["username"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
		}
	}
}

// rateLimitByToken is a middleware that rate limits according to TokenBucket and from request's JWT token
func rateLimitByToken(tb *ginratelimit.TokenBucket) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := GetBearerToken(c.GetHeader("Authorization"))
		if token != "" {
			rate := ginratelimit.RateLimitByUserId(tb, token)
			rate(c)
		}
	}
}

func main() {
	r := gin.Default()

	r.POST("/register", registerHandler)
	r.POST("/login", loginHandler)
	r.GET("/profile", rateLimitByToken(tb), authMiddleware(), getProfileHandler)
	r.PUT("/profile", rateLimitByToken(tb), authMiddleware(), putProfileHandler)
	r.GET("/admin/data", rateLimitByToken(tb), authMiddleware(), dataHandler)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
