package main

import (
	"database/sql"
	"log"
	"time"

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
	db = getDb(dsn)

	// Create a new token bucket rate limiter; ie. threshold requests per ttl time
	tb = ginratelimit.NewTokenBucket(rate_threshold, rate_ttl)
}

func main() {
	r := gin.Default()

	r.POST("/register", RegisterHandler(db))
	r.POST("/login", LoginHandler(db))
	r.GET("/profile", rateLimitByTokenMiddleware(tb), authMiddleware(db), GetProfileHandler(db))
	r.PUT("/profile", rateLimitByTokenMiddleware(tb), authMiddleware(db), PutProfileHandler(db))
	r.GET("/admin/data", rateLimitByTokenMiddleware(tb), authMiddleware(db), DataHandler(db))

	if err := r.Run(listenAddr); err != nil {
		log.Fatal(err)
	}
}
