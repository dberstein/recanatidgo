package main

import (
	"flag"
	"log"
	"time"

	"github.com/joho/godotenv"
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
var rate_threshold = 5                       // todo: parametrize
var rate_ttl = 1 * time.Minute               // todo: parametrize
var jwtSecretKey = []byte("your-secret-key") // todo: parametrize JWT secret

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addrPtr := flag.String("addr", "0.0.0.0:8080", "Listen address")

	flag.Parse()

	db := getDb(dsn)
	defer db.Close()

	// Create a new token bucket rate limiter; ie. threshold requests per ttl time
	tb := ginratelimit.NewTokenBucket(rate_threshold, rate_ttl)

	s := NewService(*addrPtr, db, tb)
	s.Serve()
}
