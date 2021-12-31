package internal

import (
	"context"
	"errors"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/donech/pearl/internal/common"
	v1 "github.com/donech/pearl/internal/proto/gen/go/gate/v1"
	systemv1 "github.com/donech/pearl/internal/proto/gen/go/system/v1"
	"github.com/donech/pearl/internal/service/system"

	"github.com/donech/pearl/internal/api/gate"
	systemApi "github.com/donech/pearl/internal/api/system"
	"github.com/donech/pearl/pkg/entry"
	grpcEntry "github.com/donech/pearl/pkg/entry/grpc"
	myJwt "github.com/donech/pearl/pkg/jwt"
	"github.com/donech/pearl/pkg/log"
	"google.golang.org/grpc"
)

func InitJwtFactory(login myJwt.LoginFunc) {
	jf, err := myJwt.NewJWTFactory(common.C().Jwt, myJwt.WithLoginFunc(login))
	if err != nil {
		log.S(context.Background()).Fatalf("init global jwt factory error: ", err)
		return
	}
	myJwt.SetJwtFactory(&jf)
}

func NewGrpcEntry(cfg grpcEntry.Config) entry.Entry {
	InitJwtFactory(Login)
	return grpcEntry.New(cfg,
		grpcEntry.WithRegisteServer(registeServer),
		grpcEntry.WithJwtFactory(myJwt.JwtFactory()),
		grpcEntry.WithJumpMethods(GetJumpMethods()),
	)
}

//GetJumpMethods 不进行 jwt 验证的 grpc handle
func GetJumpMethods() map[string]bool {
	return map[string]bool{
		"/donech.pearl.gate.v1.GateAPI/Readiness":  true,
		"/donech.pearl.gate.v1.GateAPI/Liveness":   true,
		"/donech.pearl.gate.v1.GateAPI/Login":      true,
		"/donech.pearl.gate.v1.GateAPI/SystemInfo": true,
	}
}

func registeServer(server *grpc.Server) {
	gateSrv := gate.GateAPI{}
	v1.RegisterGateAPIServer(server, gateSrv)
	systemSrv := systemApi.SystemAPI{}
	systemv1.RegisterSystemAPIServer(server, &systemSrv)
}

func Login(ctx context.Context, form myJwt.LoginForm) (jwt.MapClaims, error) {
	user, err := system.GetUserByAccount(ctx, form.Username)
	if err != nil {
		return nil, err
	}

	if !common.ValidatePassword(form.Password, user.Password) {
		return nil, errors.New("username or password error")
	}

	return jwt.MapClaims{
		"id":   strconv.FormatInt(user.ID, 10),
		"name": user.Name,
	}, nil
}
