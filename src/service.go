package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	ginratelimit "github.com/ljahier/gin-ratelimit"
)

type Service struct {
	addr string
	db   *sql.DB
	tb   *ginratelimit.TokenBucket
}

func NewService(listenAddr string, db *sql.DB, tb *ginratelimit.TokenBucket) *Service {
	return &Service{
		addr: listenAddr,
		db:   db,
		tb:   tb,
	}
}

func (s *Service) Serve() {
	r := gin.Default()

	r.POST("/register", RegisterHandler(s.db))
	r.POST("/login", LoginHandler(s.db))
	r.GET("/profile", rateLimitByTokenMiddleware(s.tb), authMiddleware(s.db), GetProfileHandler(s.db))
	r.PUT("/profile", rateLimitByTokenMiddleware(s.tb), authMiddleware(s.db), PutProfileHandler(s.db))
	r.GET("/admin/data", rateLimitByTokenMiddleware(s.tb), authMiddleware(s.db), DataHandler(s.db))

	if err := r.Run(s.addr); err != nil {
		log.Fatal(err)
	}
}
