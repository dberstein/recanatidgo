package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/dberstein/recanatid-go/models"
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

func createUser(t *testing.T, role string) (string, string, []byte) {
	assert := assert.New(t)
	router, _ := service.SetupRouter()

	w := httptest.NewRecorder()
	username := time.Now().Format("2006-01-02 15:04:05.000000000")
	exampleUser := typ.RegisterUser{
		Username: username,
		Password: "secret",
		Email:    "testing@test.com",
		Role:     role,
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
	assert.NotEqual(t, tokenRegister.(string), "")

	return username, tokenRegister.(string), userJson
}

func createAdminToken(t *testing.T, router *gin.Engine) (string, string) {
	// register user
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

	var registerResponse any
	err := json.Unmarshal(w.Body.Bytes(), &registerResponse)
	if err != nil {
		t.Fatal(err)
	}
	// token from registration
	tokenRegister := registerResponse.(map[string]any)["token"]

	// return token
	return tokenRegister.(string), username
}

func isTokenValid(token string, uri string) bool {
	router, _ := service.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	return w.Code > 100 && w.Code < 400
}

func requestWithBody(router *gin.Engine, headers map[string]string, method string, uri string, body io.Reader) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, uri, body)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	router.ServeHTTP(w, req)
	return w
}

func TestAdminDataRoute(t *testing.T) {
	assert := assert.New(t)
	router, _ := service.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin/data", nil)
	req.Header.Add("Authorization", "Bearer 123")
	router.ServeHTTP(w, req)

	assert.Equal(401, w.Code)
	assert.Equal(`{"error":"Invalid token"}`, w.Body.String())
}

func TestRegisterMissingAll(t *testing.T) {
	assert := assert.New(t)
	router, _ := service.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", nil)

	router.ServeHTTP(w, req)

	assert.Equal(400, w.Code)
	assert.Equal(`{"error":"invalid request"}`, w.Body.String())
}

func TestRegisterMissingUsername(t *testing.T) {
	assert := assert.New(t)
	router, _ := service.SetupRouter()

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
	router, _ := service.SetupRouter()

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
	router, _ := service.SetupRouter()

	w := httptest.NewRecorder()
	username := time.Now().Format("2006-01-02 15:04:05.000000000")
	exampleUser := typ.RegisterUser{
		Username: username,
		Password: "password",
	}
	userJson, _ := json.Marshal(exampleUser)
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(string(userJson)))
	router.ServeHTTP(w, req)
	assert.Equal(400, w.Code)
	assert.Equal(`{"error":"missing: email"}`, w.Body.String())
}

func TestRegister(t *testing.T) {
	assert := assert.New(t)
	router, _ := service.SetupRouter()

	_, tokenRegister, userJson := createUser(t, "admin")

	// pause
	time.Sleep(100 * time.Millisecond)

	w := httptest.NewRecorder()
	reqLogin, _ := http.NewRequest("POST", "/login", strings.NewReader(string(userJson)))
	router.ServeHTTP(w, reqLogin)
	assert.Equal(200, w.Code)

	var loginResponse any
	err := json.Unmarshal(w.Body.Bytes(), &loginResponse)
	if err != nil {
		t.Fatal(err)
	}
	// token from login
	tokenLogin := loginResponse.(map[string]any)["token"]
	assert.NotEqual(tokenLogin, "")
	// do we have a different token?
	assert.True(tokenRegister != tokenLogin)

	// test tokens
	for name, token := range map[string]string{
		"register": tokenRegister,
		"login":    tokenLogin.(string),
	} {
		assert.True(
			isTokenValid(token, "/admin/data"),
			fmt.Sprintf("token invalid: %q: %s", name, token),
		)
	}
}

func TestLogin(t *testing.T) {
	assert := assert.New(t)
	router, _ := service.SetupRouter()

	username, _, _ := createUser(t, "regular")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(fmt.Sprintf(`{"username": "%s", "password":"secret"}`, username)))
	router.ServeHTTP(w, req)
	assert.Equal(200, w.Code)

	// response body should be JSON of the form '{"token": "..."}
	var dat map[string]interface{}

	if err := json.Unmarshal(w.Body.Bytes(), &dat); err != nil {
		t.Fatal(err)
	}

	token, ok := dat["token"]
	if !ok {
		t.Fatal("response missing: token")
	}
	assert.NotEmpty(token)

	// test bad password
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/login", strings.NewReader(fmt.Sprintf(`{"username": "%s", "password":"badpassword"}`, username)))
	router.ServeHTTP(w, req)
	assert.Equal(401, w.Code)
	assert.Equal(`{"error":"Invalid credentials"}`, w.Body.String())

	// test broken json
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/login", strings.NewReader(fmt.Sprintf(`{"username": "%s", "password":`, username)))
	router.ServeHTTP(w, req)
	assert.Equal(400, w.Code)
	assert.Equal(`{"error":"unexpected EOF"}`, w.Body.String())
}

func TestGetProfileRoute(t *testing.T) {
	assert := assert.New(t)
	router, _ := service.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/profile", nil)

	router.ServeHTTP(w, req)

	assert.Equal(401, w.Code)
	assert.Equal(`{"error":"Authorization header is required"}`, w.Body.String())

	req, _ = http.NewRequest("GET", "/profile", nil)
	req.Header.Set("Authorization", "Bearer 123")
	router.ServeHTTP(w, req)

	assert.Equal(401, w.Code)
	// assert.Equal(`{"error":"Invalid token"}`, w.Body.String())

	// get token
	token, _ := createAdminToken(t, router)

	// test good request
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/profile", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	router.ServeHTTP(w, req)

	assert.Equal(200, w.Code)
}

func TestPutProfileRoute(t *testing.T) {
	assert := assert.New(t)
	router, db := service.SetupRouter()
	profile := models.NewProfile(db)

	// get token and created username
	token, username := createAdminToken(t, router)

	// test updating "role"...
	w := requestWithBody(router, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}, "PUT", "/profile", strings.NewReader(`{"role":"test"}`))
	assert.Equal(200, w.Code)

	regUser, err := profile.Get(username)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(regUser.Role, "test") // was 'admin' should bbe now 'test'

	// test update password changes pwhash...
	time.Sleep(100 * time.Millisecond)
	w = requestWithBody(router, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}, "PUT", "/profile", strings.NewReader(`{"password":"test"}`))
	assert.Equal(200, w.Code)
	regUserUpdated, err := profile.Get(username)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(regUser.Pwhash)
	assert.NotEmpty(regUserUpdated.Pwhash)
	assert.NotEqual(regUser.Pwhash, regUserUpdated.Pwhash)

	// test update email...
	w = requestWithBody(router, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}, "PUT", "/profile", strings.NewReader(`{"email":"other@other.com"}`))
	assert.Equal(200, w.Code)
	regUserUpdated, err = profile.Get(username)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(regUser.Email)
	assert.NotEmpty(regUserUpdated.Email)
	assert.NotEqual(regUser.Email, regUserUpdated.Email)
}
