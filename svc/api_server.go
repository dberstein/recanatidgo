package svc

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	ginratelimit "github.com/ljahier/gin-ratelimit"

	"github.com/dberstein/recanatid-go/handlers"
	mw "github.com/dberstein/recanatid-go/middlewares"
	"github.com/dberstein/recanatid-go/svc/owm"
	"github.com/dberstein/recanatid-go/svc/store"
	"github.com/dberstein/recanatid-go/svc/token"
)

type ApiServerOption func(*ApiServer)

type ApiServer struct {
	store    store.Storage
	tb       *ginratelimit.TokenBucket
	owmer    owm.Owmer
	jwtMaker token.JWTMaker
}

func WithStore(store store.Storage) ApiServerOption {
	return func(s *ApiServer) {
		s.store = store
	}
}

func WithTokenBucket(tb *ginratelimit.TokenBucket) ApiServerOption {
	return func(s *ApiServer) {
		s.tb = tb
	}
}

func WithOwmer(o owm.Owmer) ApiServerOption {
	return func(s *ApiServer) {
		s.owmer = o
	}
}

func WithJWTMaker(jwtMaker token.JWTMaker) ApiServerOption {
	return func(s *ApiServer) {
		s.jwtMaker = jwtMaker
	}
}

func NewApiServer(option ...ApiServerOption) *ApiServer {
	s := &ApiServer{}
	for _, o := range option {
		o(s)
	}
	return s
}

func (s *ApiServer) SetupRouter() (*gin.Engine, *sql.DB) {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	db := s.store.GetDB()

	r.POST("/register", handlers.RegisterHandler(db, s.jwtMaker))
	r.POST("/login", handlers.LoginHandler(db, s.jwtMaker))
	r.GET("/profile", mw.RateLimitByTokenMiddleware(s.tb), mw.AuthMiddleware(s.jwtMaker), handlers.GetProfileHandler(db))
	r.PUT("/profile", mw.RateLimitByTokenMiddleware(s.tb), mw.AuthMiddleware(s.jwtMaker), handlers.PutProfileHandler(db))
	r.GET("/admin/data", mw.RateLimitByTokenMiddleware(s.tb), mw.AuthMiddleware(s.jwtMaker), mw.RoleMiddleware([]string{"admin"}), handlers.DataHandler(db, s.owmer))

	return r, db
}

func (s *ApiServer) Serve(addr string) error {
	r, _ := s.SetupRouter()
	srv := &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Serving:", addr)
	return srv.ListenAndServe()
}
