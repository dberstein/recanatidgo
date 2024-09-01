package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/dberstein/recanatid-go/svc"
	"github.com/dberstein/recanatid-go/svc/db"
	"github.com/dberstein/recanatid-go/svc/owm"
	"github.com/dberstein/recanatid-go/svc/rate"
	"github.com/dberstein/recanatid-go/svc/store"
	"github.com/dberstein/recanatid-go/svc/token"

	_ "github.com/mattn/go-sqlite3"
)

var dsn string
var jwtSecretKey = []byte("your-secret-key") // todo: parametrize JWT secret

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

	rl, err := rate.NewRateLimiter(*ratePtr)
	if err != nil {
		log.Fatal(err)
	}

	s := svc.NewApiServer(
		svc.WithStore(store.NewStore(dbcon)),
		svc.WithTokenBucket(rl.GetTokenBucket()),
		svc.WithOwmer(owm.NewOwm(*owmPtr)),
		svc.WithJMWMaker(token.NewJWTMaker(jwtSecretKey)),
	)

	if err := s.Serve(*addrPtr); err != nil {
		log.Fatal(err)
	}
}
