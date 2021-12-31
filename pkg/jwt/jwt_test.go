package jwt

import (
	"context"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

var factory JWTFactory

func init() {
	f, _ := NewJWTFactory(Config{
		SingingAlgorithm: "RS256",
		Key:              "asdasdada123123",
		PublicKeyFile:    "./testData/zq_mall_rsa_public_key.pem",
		PrivateKeyFile:   "./testData/zq_mall_rsa_private_key.pem",
		Timeout:          "10000m",
	}, WithLoginFunc(login))
	factory = f
}

func TestJWTFactory_GenerateToken(t *testing.T) {
	ctx := context.Background()
	token, err := factory.GenerateToken(ctx, LoginForm{
		Username: "12312",
		Password: "123123",
	})
	if err != nil {
		t.Fatal("generate token failed: ", err.Error())
	}
	t.Log("generate token success and token is: ", token)
}

func TestJWTFactory_VerifyToken(t *testing.T) {
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTY3OTYzOTEsImlkIjoxLCJpZGVudGlmeSI6ImVycm9ycyIsInVzZXJuYW1lIjoiMTIzMTIifQ.XNrAPJnWQOBaZhoDfvlOqxGqubPqOVKlaXlURwDg2ng63ZWwPJ4ZQ-8vIPUK9StF9UT-MXTNep1TtASWBUt58bGIZ4mVvkce9owjqNgZDBamS4650KU6nJJ6l_PgG5MHymuu5Ehj5xJcz_LIZPl-iXPc7T_Xh1b-xJLqQFb13MCyPWxvFymHGTEZ6a4mf8svVwLWBmdyglnddVUWUnkZRD1nQILXpvOtlHUhUOWFbnEAmLDB8MKy6CCv3C_rBfk0ZgUYwWVP2hqWNVXuhc05UIOkKTUnYF-TEV0aByeFYRNl_VxJJ9Ix3Gm7ZmihzY5RzYhZnVjPeBp4XaClS_wlXQ"
	if !factory.VerifyToken(token) {
		t.Fatal("verify token failed")
	}
	t.Log("verify token success")
}

func TestJWTFactory_GetPayload(t *testing.T) {
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTY3OTYzOTEsImlkIjoxLCJpZGVudGlmeSI6ImVycm9ycyIsInVzZXJuYW1lIjoiMTIzMTIifQ.XNrAPJnWQOBaZhoDfvlOqxGqubPqOVKlaXlURwDg2ng63ZWwPJ4ZQ-8vIPUK9StF9UT-MXTNep1TtASWBUt58bGIZ4mVvkce9owjqNgZDBamS4650KU6nJJ6l_PgG5MHymuu5Ehj5xJcz_LIZPl-iXPc7T_Xh1b-xJLqQFb13MCyPWxvFymHGTEZ6a4mf8svVwLWBmdyglnddVUWUnkZRD1nQILXpvOtlHUhUOWFbnEAmLDB8MKy6CCv3C_rBfk0ZgUYwWVP2hqWNVXuhc05UIOkKTUnYF-TEV0aByeFYRNl_VxJJ9Ix3Gm7ZmihzY5RzYhZnVjPeBp4XaClS_wlXQ"
	claims, err := factory.GetClaims(token)
	if err != nil {
		t.Fatal("get payload failed: ", err.Error())
	}
	t.Log("get payload success: ", claims)
}

func login(ctx context.Context, form LoginForm) (jwt.MapClaims, error) {
	return jwt.MapClaims{"username": form.Username, "id": 1, "identify": "errors"}, nil
}
