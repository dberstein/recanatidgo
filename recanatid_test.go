package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dberstein/recanatid-go/svc"
	"github.com/dberstein/recanatid-go/svc/db"
	"github.com/dberstein/recanatid-go/svc/rate"
	"github.com/dberstein/recanatid-go/svc/store"
	"github.com/dberstein/recanatid-go/svc/token"
)

var service *svc.ApiServer

func init() {
	rl, err := rate.NewRateLimiter("1/1m")
	if err != nil {
		log.Fatal(err)
	}

	dbcon := db.NewDb(":memory:")

	service = svc.NewApiServer(
		svc.WithStore(store.NewStore(dbcon)),
		svc.WithTokenBucket(rl.GetTokenBucket()),
		// svc.WithOwmer(owm.NewOwm(*owmPtr)),
		svc.WithJWTMaker(token.NewJWTMaker([]byte("secret"))),
	)
}

func TestAdminDataRoute(t *testing.T) {
	router := service.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin/data", nil)
	req.Header.Add("Authorization", "Bearer 123")
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, `{"error":"Invalid token"}`, w.Body.String())
}
