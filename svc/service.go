package svc

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	ginratelimit "github.com/ljahier/gin-ratelimit"

	"github.com/dberstein/recanatid-go/handler"
	mw "github.com/dberstein/recanatid-go/middleware"
	"github.com/dberstein/recanatid-go/svc/owm"
	"github.com/dberstein/recanatid-go/svc/token"
)

type Service struct {
	db       *sql.DB
	tb       *ginratelimit.TokenBucket
	owmer    owm.Owmer
	jwtMaker *token.JWTMaker
}

func NewService(db *sql.DB, tb *ginratelimit.TokenBucket, owmer owm.Owmer, jwtMaker *token.JWTMaker) *Service {
	return &Service{
		db:       db,
		tb:       tb,
		owmer:    owmer,
		jwtMaker: jwtMaker,
	}
}

func (s *Service) Serve(addr string) error {
	r := gin.Default()

	r.POST("/register", handler.RegisterHandler(s.db, s.jwtMaker))
	r.POST("/login", handler.LoginHandler(s.db, s.jwtMaker))
	r.GET("/profile", mw.RateLimitByTokenMiddleware(s.tb), mw.AuthMiddleware(s.jwtMaker), handler.GetProfileHandler(s.db))
	r.PUT("/profile", mw.RateLimitByTokenMiddleware(s.tb), mw.AuthMiddleware(s.jwtMaker), handler.PutProfileHandler(s.db))
	r.GET("/admin/data", mw.RateLimitByTokenMiddleware(s.tb), mw.AuthMiddleware(s.jwtMaker), handler.DataHandler(s.owmer))

	return r.Run(addr)
}
