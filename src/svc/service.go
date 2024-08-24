package svc

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	ginratelimit "github.com/ljahier/gin-ratelimit"

	"github.com/dberstein/recanatid-go/src/handler"
	"github.com/dberstein/recanatid-go/src/mw"
	"github.com/dberstein/recanatid-go/src/owm"
	"github.com/dberstein/recanatid-go/src/token"
)

type Service struct {
	addr     string
	db       *sql.DB
	tb       *ginratelimit.TokenBucket
	owmer    owm.Owmer
	jwtMaker *token.JWTMaker
}

func NewService(addr string, db *sql.DB, tb *ginratelimit.TokenBucket, owmer owm.Owmer, jwtMaker *token.JWTMaker) *Service {
	return &Service{
		addr:     addr,
		db:       db,
		tb:       tb,
		owmer:    owmer,
		jwtMaker: jwtMaker,
	}
}

func (s *Service) Serve() {
	r := gin.Default()

	r.POST("/register", handler.RegisterHandler(s.db, s.jwtMaker))
	r.POST("/login", handler.LoginHandler(s.db, s.jwtMaker))
	r.GET("/profile", mw.RateLimitByTokenMiddleware(s.tb), mw.AuthMiddleware(s.jwtMaker), handler.GetProfileHandler(s.db))
	r.PUT("/profile", mw.RateLimitByTokenMiddleware(s.tb), mw.AuthMiddleware(s.jwtMaker), handler.PutProfileHandler(s.db))
	r.GET("/admin/data", mw.RateLimitByTokenMiddleware(s.tb), mw.AuthMiddleware(s.jwtMaker), handler.DataHandler(s.owmer))

	if err := r.Run(s.addr); err != nil {
		log.Fatal(err)
	}
}
