package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	ginratelimit "github.com/ljahier/gin-ratelimit"

	"github.com/dberstein/recanatid-go/svc"
	"github.com/dberstein/recanatid-go/svc/db"
	"github.com/dberstein/recanatid-go/svc/owm"
	"github.com/dberstein/recanatid-go/svc/token"

	_ "github.com/mattn/go-sqlite3"
)

var dsn string
var jwtSecretKey = []byte("your-secret-key") // todo: parametrize JWT secret

func parseRateString(s string) (int, time.Duration, error) {
	rateTmp := strings.Split(s, "/")
	rate, err := strconv.Atoi(rateTmp[0])
	if err != nil {
		return 0, 0, err
	}
	if rate <= 0 {
		return rate, 0, errors.New("rate has to be greater than zero")
	}

	ttl, err := time.ParseDuration(rateTmp[1])
	if err != nil {
		return rate, 0, err
	}
	if ttl <= 0 {
		return rate, ttl, errors.New("rate's TTL has to be greater than zero")
	}
	return rate, ttl, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addrPtr := flag.String("addr", "0.0.0.0:8080", "Listen address")
	dsnPtr := flag.String("dsn", "x.db", "Database DSN")
	owmPtr := flag.String("owm", os.Getenv("OWM_API_KEY"), "OWM API key")
	ratePtr := flag.String("rate", "5/1m", "Rate limit string: \"<rate>/<ttl>\"")

	flag.Parse()

	dsn = *dsnPtr

	dbcon := db.NewDb(dsn)
	defer dbcon.Close()

	rate, ttl, err := parseRateString(*ratePtr)
	if err != nil {
		log.Fatal(err)
	}

	s := svc.NewApiService(
		svc.WithDB(dbcon),
		svc.WithTokenBucket(ginratelimit.NewTokenBucket(rate, ttl)),
		svc.WithOwmer(owm.NewOwm(*owmPtr)),
		svc.WithJMWMaker(token.NewJWTMaker(jwtSecretKey)),
	)

	if err := s.Serve(*addrPtr); err != nil {
		log.Fatal(err)
	}
}
