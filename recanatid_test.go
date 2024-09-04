package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dberstein/recanatid-go/svc"
	"github.com/dberstein/recanatid-go/svc/db"
	"github.com/dberstein/recanatid-go/svc/rate"
	"github.com/dberstein/recanatid-go/svc/store"
	"github.com/dberstein/recanatid-go/svc/token"
	"github.com/dberstein/recanatid-go/typ"
)

var service *svc.ApiServer

func init() {
	rl, err := rate.NewRateLimiter("1000/1m")
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

func TestRegisterMissingAll(t *testing.T) {
	router := service.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"error":"invalid request"}`, w.Body.String())
}

func TestRegisterMissingUsername(t *testing.T) {
	assert := assert.New(t)
	router := service.SetupRouter()

	w := httptest.NewRecorder()

	// missing username
	exampleUser := typ.RegisterUser{
		Password: "password",
		Email:    "testing@test.com",
	}
	userJson, _ := json.Marshal(exampleUser)
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(string(userJson)))
	router.ServeHTTP(w, req)
	assert.Equal(400, w.Code)
	assert.Equal(`{"error":"missing: username"}`, w.Body.String())
}

func TestRegisterMissingPassword(t *testing.T) {
	assert := assert.New(t)
	router := service.SetupRouter()

	w := httptest.NewRecorder()
	username := time.Now().Format("2006-01-02 15:04:05.000000000")
	exampleUser := typ.RegisterUser{
		Username: username,
		Email:    "testing@test.com",
	}
	userJson, _ := json.Marshal(exampleUser)
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(string(userJson)))
	router.ServeHTTP(w, req)
	assert.Equal(400, w.Code)
	assert.Equal(`{"error":"missing: password"}`, w.Body.String())
}

func TestRegisterMissingEmail(t *testing.T) {
	assert := assert.New(t)
	router := service.SetupRouter()

	w := httptest.NewRecorder()
	username := time.Now().Format("2006-01-02 15:04:05.000000000")
	exampleUser := typ.RegisterUser{
		Username: username,
		Email:    "testing@test.com",
	}
	userJson, _ := json.Marshal(exampleUser)
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(string(userJson)))
	router.ServeHTTP(w, req)
	assert.Equal(400, w.Code)
	assert.Equal(`{"error":"missing: password"}`, w.Body.String())
}

func TestRegister(t *testing.T) {
	assert := assert.New(t)
	router := service.SetupRouter()

	w := httptest.NewRecorder()
	username := time.Now().Format("2006-01-02 15:04:05.000000000")
	exampleUser := typ.RegisterUser{
		Username: username,
		Password: "secret",
		Email:    "testing@test.com",
		Role:     "admin",
	}
	userJson, _ := json.Marshal(exampleUser)
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(string(userJson)))
	router.ServeHTTP(w, req)
	assert.Equal(200, w.Code)

	var registerResponse any
	err := json.Unmarshal(w.Body.Bytes(), &registerResponse)
	if err != nil {
		t.Fatal(err)
	}
	// token from registration
	tokenRegister := registerResponse.(map[string]any)["token"]
	assert.NotEqual(t, tokenRegister, "")

	w = httptest.NewRecorder()
	reqLogin, _ := http.NewRequest("POST", "/login", strings.NewReader(string(userJson)))
	router.ServeHTTP(w, reqLogin)
	assert.Equal(200, w.Code)

	var loginResponse any
	err = json.Unmarshal(w.Body.Bytes(), &loginResponse)
	if err != nil {
		t.Fatal(err)
	}
	// token from login
	tokenLogin := loginResponse.(map[string]any)["token"]
	assert.NotEqual(tokenLogin, "")
	// do we have a different token?
	assert.True(tokenRegister != tokenLogin)

	// // test tokens
	// for name, token := range map[string]string{"register": tokenRegister.(string), "login": tokenLogin.(string)} {
	// 	assert.True(t, validToken(t, token), fmt.Sprintf("token invalid: %q: %s", name, token))
	// }
}

func validToken(t *testing.T, token string) bool {
	router := service.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin/data", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	t.Log("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	t.Log("Code", w.Code)

	return w.Code > 100 && w.Code < 400
	// assert.Equal(t, 401, w.Code)
	// assert.Equal(t, `{"error":"Invalid token"}`, w.Body.String())

	// return false
}

func TestLogin(t *testing.T) {
	assert := assert.New(t)
	router := service.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login", nil)

	router.ServeHTTP(w, req)

	assert.Equal(404, w.Code)
	assert.Equal(`404 page not found`, w.Body.String())
}

func TestGetProfileRoute(t *testing.T) {
	assert := assert.New(t)
	router := service.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/profile", nil)

	router.ServeHTTP(w, req)

	assert.Equal(401, w.Code)
	assert.Equal(`{"error":"Authorization header is required"}`, w.Body.String())

	req, _ = http.NewRequest("GET", "/profile", nil)
	req.Header.Add("Authorizationx", "Bearer 123")
	assert.Equal(401, w.Code)
	assert.Equal(`{"error":"Authorization header is required"}`, w.Body.String())
}
