package gate

import (
	"context"
	"fmt"
	"time"

	"github.com/donech/pearl/internal/common"
	v1 "github.com/donech/pearl/internal/proto/gen/go/gate/v1"
	"github.com/donech/pearl/pkg/jwt"
	"github.com/donech/pearl/pkg/log"
)

type GateAPI struct{}

//Login.
func (i GateAPI) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	token, err := jwt.JwtFactory().GenerateToken(ctx, jwt.LoginForm{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		log.S(ctx).Warn("Login error, ", err)
		return &v1.LoginResponse{
			Code: common.ErrorCode,
			Msg:  common.ResponseMsg(common.ErrorCode),
		}, nil
	}
	return &v1.LoginResponse{
		Code:  common.SuccessCode,
		Msg:   common.ResponseMsg(common.SuccessCode),
		Token: token,
	}, err
}

//Readiness.
func (i GateAPI) Readiness(_ context.Context, _ *v1.ReadinessRequest) (*v1.ReadinessResponse, error) {
	return &v1.ReadinessResponse{
		Code: 0,
		Msg:  "pong",
	}, nil
}

//Liveness.
func (i GateAPI) Liveness(_ context.Context, _ *v1.LivenessRequest) (*v1.LivenessResponse, error) {
	return &v1.LivenessResponse{
		Code: 0,
		Msg:  "pong",
	}, nil
}

//SystemInfo.
func (i GateAPI) SystemInfo(_ context.Context, _ *v1.SystemInfoRequest) (*v1.SystemInfoResponse, error) {
	c := common.C()
	return &v1.SystemInfoResponse{
		Code: common.SuccessCode,
		Msg:  common.ResponseMsg(common.SuccessCode),
		Data: &v1.SystemInfoResponse_Data{
			Name:      c.SystemInfo.Name,
			Version:   c.SystemInfo.Version,
			Copyright: fmt.Sprintf("%d %s", time.Now().Year(), c.SystemInfo.Copyright),
			Desc:      c.SystemInfo.Desc,
			Href:      c.SystemInfo.Href,
		},
	}, nil
}
