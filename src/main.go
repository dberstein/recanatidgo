package main

import (
	"time"

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

var dsn = "x.db"                             // todo: parametrize
var listenAddr string = ":8080"              // todo: parametrize
var rate_threshold = 5                       // todo: parametrize
var rate_ttl = 1 * time.Minute               // todo: parametrize
var jwtSecretKey = []byte("your-secret-key") // todo: parametrize JWT secret

func main() {
	db := getDb(dsn)
	defer db.Close()

	// Create a new token bucket rate limiter; ie. threshold requests per ttl time
	tb := ginratelimit.NewTokenBucket(rate_threshold, rate_ttl)

	s := NewService(listenAddr, db, tb)
	s.Serve()
}
