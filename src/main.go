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

// See init() for initialization
var db *sql.DB
var tb *ginratelimit.TokenBucket

var dsn = "x.db"                             // todo: parametrize
var listenAddr string = ":8080"              // todo: parametrize
var rate_threshold = 5                       // todo: parametrize
var rate_ttl = 1 * time.Minute               // todo: parametrize
var jwtSecretKey = []byte("your-secret-key") // todo: parametrize JWT secret

func init() {
	db = GetDb()

	// Create a new token bucket rate limiter; ie. threshold requests per ttl time
	tb = ginratelimit.NewTokenBucket(rate_threshold, rate_ttl)
}

func ensureSchema(db *sql.DB) error {
	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS users (username TEXT PRIMARY KEY, pwhash TEXT, email TEXT, role TEXT)",
	); err != nil {
		return err
	}
	return nil
}

func GetDb() *sql.DB {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = ensureSchema(db)
	if err != nil {
		log.Fatal(err)
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

	r.POST("/register", RegisterHandler)
	r.POST("/login", LoginHandler)
	r.GET("/profile", rateLimitByToken(tb), authMiddleware(), GetProfileHandler)
	r.PUT("/profile", rateLimitByToken(tb), authMiddleware(), PutProfileHandler)
	r.GET("/admin/data", rateLimitByToken(tb), authMiddleware(), DataHandler)

	if err := r.Run(listenAddr); err != nil {
		log.Fatal(err)
	}
}
