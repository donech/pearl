package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/donech/pearl/pkg/jwt"
	"github.com/donech/pearl/pkg/log"
	"github.com/gin-gonic/gin"
)

var TokenNotFoundErr = errors.New("token not found err")

var errCode = 1
var tokenName = "token"

type JWTMiddleware struct {
	factory            jwt.JWTFactory
	tokenLookup        string
	responseHandler    ResponseHandler
	errResponseHandler ResponseHandler
	successCode        int
	errorCode          int
	tokenName          string
}

type ResponseHandler func(ctx *gin.Context, code int, msg string, data interface{})

func DefaultResponseHandler(ctx *gin.Context, code int, msg string, data interface{}) {
	ctx.JSON(200, gin.H{"code": code, "msg": msg, "data": data})
}

type Option func(middleware *JWTMiddleware)

func WithFactory(factory jwt.JWTFactory) Option {
	return func(middleware *JWTMiddleware) {
		middleware.factory = factory
	}
}

func WithTokenLookup(tokenLookup string) Option {
	return func(middleware *JWTMiddleware) {
		middleware.tokenLookup = tokenLookup
	}
}

func WithSuccessCode(code int) Option {
	return func(m *JWTMiddleware) {
		m.successCode = code
	}
}

func WithErrorCode(code int) Option {
	return func(m *JWTMiddleware) {
		m.errorCode = code
	}
}

func WithTokenName(name string) Option {
	return func(m *JWTMiddleware) {
		m.tokenName = name
	}
}

func WithResponseHandler(handler ResponseHandler) Option {
	return func(m *JWTMiddleware) {
		m.responseHandler = handler
	}
}

func WithErrResponseHandler(handler ResponseHandler) Option {
	return func(m *JWTMiddleware) {
		m.errResponseHandler = handler
	}
}

func NewJWTMiddleware(opts ...Option) JWTMiddleware {
	m := JWTMiddleware{}
	for _, o := range opts {
		o(&m)
	}
	if m.responseHandler == nil {
		m.responseHandler = DefaultResponseHandler
	}
	if m.errResponseHandler == nil {
		m.errResponseHandler = DefaultResponseHandler
	}

	if m.errorCode == 0 {
		m.errorCode = errCode
	}

	if m.tokenName == "" {
		m.tokenName = tokenName
	}

	if m.tokenLookup == "" {
		m.tokenLookup = "header:" + m.tokenName
	}
	return m
}

func (j JWTMiddleware) GenerateTokenHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		f := jwt.LoginForm{}
		err := ctx.ShouldBind(&f)
		if err != nil {
			j.errResponseHandler(ctx, j.errorCode, err.Error(), nil)
			return
		}
		token, err := j.factory.GenerateToken(ctx.Request.Context(), f)
		if err != nil {
			j.errResponseHandler(ctx, j.errorCode, err.Error(), nil)
			return
		}
		j.responseHandler(ctx, j.successCode, "success", gin.H{"token": token})
	}
}

func (j JWTMiddleware) MiddleWareImpl() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := j.GetToken(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"msg": "get jwt-token failed: " + err.Error()})
			return
		}
		claims, err := j.factory.GetClaims(token)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"msg": "get jwt-claims failed: " + err.Error()})
			return
		}
		c := context.WithValue(ctx.Request.Context(), jwt.CtxJWTKey, claims)
		ctx.Request = ctx.Request.WithContext(c)
		log.S(ctx.Request.Context()).Infof("set claims to ctx, key=%#v, claims=%#v", jwt.CtxJWTKey, claims)
		ctx.Next()
	}
}

func (j JWTMiddleware) GetToken(ctx *gin.Context) (token string, err error) {
	scopes := strings.Split(j.tokenLookup, ",")
	for _, scope := range scopes {
		if len(token) > 0 {
			break
		}
		parts := strings.Split(scope, ":")
		if len(parts) != 2 {
			continue
		}
		method, key := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		switch method {
		case "header":
			token, err = getTokenFromHeader(ctx, key)
		}
	}
	if token == "" {
		err = TokenNotFoundErr
	}
	return token, err
}

func getTokenFromHeader(ctx *gin.Context, key string) (token string, err error) {
	token = ctx.Request.Header.Get(key)
	return
}
