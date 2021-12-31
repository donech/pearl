package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	myJwt "github.com/donech/pearl/pkg/jwt"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

var factory = newJWTFactory()

func newJWTFactory() myJwt.JWTFactory {
	f, err := myJwt.NewJWTFactory(
		myJwt.Config{
			SingingAlgorithm: "HS256",
			Key:              "12312asdasd",
			PublicKeyFile:    "",
			PrivateKeyFile:   "",
			Timeout:          "10m",
		},
		myJwt.WithLoginFunc(login),
	)
	if err != nil {
		log.Fatal("")
	}
	return f
}

func TestGenerateTokenHandler(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/getToken", bytes.NewBufferString(`{"username":"1111","password":"2222"}`))
	c.Request.Header.Add("Content-Type", gin.MIMEJSON)
	middleware := NewJWTMiddleware(WithFactory(factory), WithTokenLookup("header:auth"))
	f := middleware.GenerateTokenHandler()
	f(c)
	resp := gin.H{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Error("json.Unmarshal failed: ", err.Error())
	}
	if code := resp["code"]; code != float64(0) {
		t.Fatal("get token api return fail, result is: ", w.Body.String())
	}
	t.Log("get token api return success, result is: ", w.Body.String())
}

func TestJWTMiddleware_MiddleWareImpl(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/getToken", bytes.NewBufferString(`{"username":"1111","password":"2222"}`))
	c.Request.Header.Add("auth", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTY0MTY4OTQsImlkIjoxLCJpZGVudGlmeSI6ImVycm9ycyIsInVzZXJuYW1lIjoiMTExMSJ9.7sGAYcSJqmZP-P-08N-GxiPGkBYlPk69KnRC0S-H4D4")
	middleware := NewJWTMiddleware(WithFactory(factory), WithTokenLookup("header:auth"))
	f := middleware.MiddleWareImpl()
	f(c)
	if w.Code == 401 {
		log.Fatal("token is not valid ", w.Body.String())
	}
	log.Println(c.Get("jwt"))
}

func login(ctx context.Context, form myJwt.LoginForm) (jwt.MapClaims, error) {
	return jwt.MapClaims{"username": form.Username, "id": 1, "identify": "errors"}, nil
}
