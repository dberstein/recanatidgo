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

func NewService(addr string, db *sql.DB, tb *ginratelimit.TokenBucket) *Service {
	return &Service{addr: addr, db: db, tb: tb}
}

func (s *Service) Serve(owm *Owm) {
	r := gin.Default()

	r.POST("/register", registerHandler(s.db))
	r.POST("/login", loginHandler(s.db))
	r.GET("/profile", rateLimitByTokenMiddleware(s.tb), authMiddleware(), getProfileHandler(s.db))
	r.PUT("/profile", rateLimitByTokenMiddleware(s.tb), authMiddleware(), putProfileHandler(s.db))
	r.GET("/admin/data", rateLimitByTokenMiddleware(s.tb), authMiddleware(), dataHandler(owm))

	if err := r.Run(s.addr); err != nil {
		log.Fatal(err)
	}
}
