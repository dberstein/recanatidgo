package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	ginratelimit "github.com/ljahier/gin-ratelimit"

	"github.com/dberstein/recanatid-go/src/owm"
	"github.com/dberstein/recanatid-go/src/token"
)

type Service struct {
	addr     string
	db       *sql.DB
	tb       *ginratelimit.TokenBucket
	owm      owm.Owmer
	jwtMaker *token.JWTMaker
}

func NewService(addr string, db *sql.DB, tb *ginratelimit.TokenBucket, owm owm.Owmer, jwtMaker *token.JWTMaker) *Service {
	return &Service{
		addr:     addr,
		db:       db,
		tb:       tb,
		owm:      owm,
		jwtMaker: jwtMaker,
	}
}

func (s *Service) Serve() {
	r := gin.Default()

	r.POST("/register", registerHandler(s.db, s.jwtMaker))
	r.POST("/login", loginHandler(s.db, s.jwtMaker))
	r.GET("/profile", rateLimitByTokenMiddleware(s.tb), authMiddleware(s.jwtMaker), getProfileHandler(s.db))
	r.PUT("/profile", rateLimitByTokenMiddleware(s.tb), authMiddleware(s.jwtMaker), putProfileHandler(s.db))
	r.GET("/admin/data", rateLimitByTokenMiddleware(s.tb), authMiddleware(s.jwtMaker), dataHandler(s.owm))

	if err := r.Run(s.addr); err != nil {
		log.Fatal(err)
	}
}
