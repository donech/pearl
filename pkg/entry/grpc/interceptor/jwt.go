package interceptor

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/donech/pearl/pkg/log"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"

	"github.com/donech/pearl/pkg/jwt"

	"google.golang.org/grpc"
)

var (
	headerAuthorize = "authorization"
)

type JwtInterceptor struct {
	jwtFactory  *jwt.JWTFactory
	jumpMethods map[string]bool
}

func NewJwtInterceptor(jwtFactory *jwt.JWTFactory, jumps map[string]bool) *JwtInterceptor {
	return &JwtInterceptor{jwtFactory: jwtFactory, jumpMethods: jumps}
}

func (i *JwtInterceptor) Serve(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if i.jumpMethods == nil {
		return handler(ctx, req)
	}
	if i.jumpMethods[info.FullMethod] {
		return handler(ctx, req)
	}
	token := metautils.ExtractIncoming(ctx).Get(headerAuthorize)
	if token == "" {
		log.S(ctx).Error("no token found")
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}
	claims, err := i.jwtFactory.GetClaims(token)
	if err != nil {
		log.S(ctx).Error("jwt GetClaims error, %+v", err)
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}
	newCtx := jwt.SetClaimsToCtx(ctx, claims)
	return handler(newCtx, req)

}

func Jwt(jwtF *jwt.JWTFactory, jumps map[string]bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if jumps == nil {
			return handler(ctx, req)
		}
		if jumps[info.FullMethod] {
			return handler(ctx, req)
		}
		token := metautils.ExtractIncoming(ctx).Get(headerAuthorize)
		if token == "" {
			log.S(ctx).Error("no token found")
			return nil, status.Error(codes.Unauthenticated, "Unauthorized")
		}
		claims, err := jwtF.GetClaims(token)
		if err != nil {
			log.S(ctx).Error("jwt GetClaims error, %+v", err)
			return nil, status.Error(codes.Unauthenticated, "Unauthorized")
		}
		newCtx := jwt.SetClaimsToCtx(ctx, claims)
		return handler(newCtx, req)
	}
}
