package svc

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	ginratelimit "github.com/ljahier/gin-ratelimit"

	"github.com/dberstein/recanatid-go/handler"
	mw "github.com/dberstein/recanatid-go/middleware"
	"github.com/dberstein/recanatid-go/svc/owm"
	"github.com/dberstein/recanatid-go/svc/token"
)

type Option func(*Service)

type Service struct {
	db       *sql.DB
	tb       *ginratelimit.TokenBucket
	owmer    owm.Owmer
	jwtMaker *token.JWTMaker
}

func WithDB(db *sql.DB) Option {
	return func(s *Service) {
		s.db = db
	}
}

func WithTokenBucket(tb *ginratelimit.TokenBucket) Option {
	return func(s *Service) {
		s.tb = tb
	}
}

func WithOwmer(o owm.Owmer) Option {
	return func(s *Service) {
		s.owmer = o
	}
}

func WithJMWMaker(jwtMaker *token.JWTMaker) Option {
	return func(s *Service) {
		s.jwtMaker = jwtMaker
	}
}

func NewService(option ...Option) *Service {
	s := &Service{}
	for _, o := range option {
		o(s)
	}
	return s
}

func (s *Service) Serve(addr string) error {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.POST("/register", handler.RegisterHandler(s.db, s.jwtMaker))
	r.POST("/login", handler.LoginHandler(s.db, s.jwtMaker))
	r.GET("/profile", mw.RateLimitByTokenMiddleware(s.tb), mw.AuthMiddleware(s.jwtMaker), handler.GetProfileHandler(s.db))
	r.PUT("/profile", mw.RateLimitByTokenMiddleware(s.tb), mw.AuthMiddleware(s.jwtMaker), handler.PutProfileHandler(s.db))
	r.GET("/admin/data", mw.RateLimitByTokenMiddleware(s.tb), mw.AuthMiddleware(s.jwtMaker), handler.DataHandler(s.db, s.owmer))

	log.Println("Serving:", addr)
	return r.Run(addr)
}
