package main

import (
	"flag"
	"log"
	"os"
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

var dsn string
var rate_threshold = 5                       // todo: parametrize
var rate_ttl = 1 * time.Minute               // todo: parametrize
var jwtSecretKey = []byte("your-secret-key") // todo: parametrize JWT secret

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addrPtr := flag.String("addr", "0.0.0.0:8080", "Listen address")
	dsnPtr := flag.String("dsn", "x.db", "Database DSN")
	owmPtr := flag.String("owm", os.Getenv("OWM_API_KEY"), "OWM API key")

	flag.Parse()

	dsn = *dsnPtr

	db := getDb(dsn)
	defer db.Close()

	s := NewService(
		*addrPtr,
		db,
		ginratelimit.NewTokenBucket(rate_threshold, rate_ttl),
		NewOwm(*owmPtr),
	)
	s.Serve()
}
