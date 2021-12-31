package system

import (
	"context"

	"github.com/donech/pearl/internal/common"
	v1 "github.com/donech/pearl/internal/proto/gen/go/system/v1"
	"github.com/donech/pearl/internal/service/system"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SystemAPI struct{}

//Users.
func (s *SystemAPI) Users(ctx context.Context, req *v1.UsersRequest) (*v1.UsersResponse, error) {
	params := system.UsersParams{
		Name:    req.GetName(),
		Account: req.GetAccount(),
	}
	if req.GetPager() != nil {
		params.Page = int(req.Pager.Page)
		params.Size = int(req.Pager.Size)
	}
	users, err := system.Users(ctx, params)
	if err != nil {
		return nil, err
	}
	respUsers := make([]*v1.User, len(users.Users))
	for i, v := range users.Users {
		respUsers[i] = &v1.User{
			Id:          v.ID,
			Name:        v.Name,
			Account:     v.Account,
			CreatedTime: v.CreatedTime.Local().String(),
			UpdatedTime: v.UpdatedTime.Local().String(),
			DeletedTime: v.DeletedTime.Local().String(),
		}
	}
	return &v1.UsersResponse{
		Code: common.SuccessCode,
		Msg:  common.ResponseMsg(common.SuccessCode),
		Data: &v1.UsersResponse_Data{
			Total: int32(users.Total),
			List:  respUsers,
		},
	}, nil
}

//User 用户详情.
func (s *SystemAPI) User(ctx context.Context, req *v1.UserRequest) (*v1.UserResponse, error) {
	user, err := system.GetUserById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.UserResponse{
		Code: common.SuccessCode,
		Msg:  common.ResponseMsg(common.SuccessCode),
		Data: &v1.User{
			Id:          user.ID,
			Name:        user.Name,
			Account:     user.Account,
			CreatedTime: user.CreatedTime.Local().String(),
			UpdatedTime: user.UpdatedTime.Local().String(),
			DeletedTime: user.DeletedTime.Local().String(),
		},
	}, nil
}

//CurrentUser 当前用户详情.
func (s *SystemAPI) CurrentUser(ctx context.Context, req *v1.CurrentUserRequest) (*v1.CurrentUserResponse, error) {
	user, err := common.AuthUser(ctx, true)
	if err != nil {
		return nil, nil
	}
	return &v1.CurrentUserResponse{
		Code: common.SuccessCode,
		Msg:  common.ResponseMsg(common.SuccessCode),
		Data: &v1.User{
			Id:          user.ID,
			Name:        user.Name,
			Account:     user.Account,
			CreatedTime: user.CreatedTime.Local().String(),
			UpdatedTime: user.UpdatedTime.Local().String(),
			DeletedTime: user.DeletedTime.Local().String(),
		},
	}, nil
}

//SaveUser 创建.
func (s *SystemAPI) SaveUser(ctx context.Context, req *v1.SaveUserRequest) (*v1.SaveUserResponse, error) {
	if req.Password != req.Confirm {
		return nil, status.Error(codes.InvalidArgument, "密码不一致!")
	}
	user, err := system.SaveUser(ctx, req.Account, req.Name, common.EncryptPassword(req.Password))
	if err != nil {
		return nil, err
	}
	return &v1.SaveUserResponse{
		Code: common.SuccessCode,
		Msg:  common.ResponseMsg(common.SuccessCode),
		Data: &v1.User{
			Id:          user.ID,
			Name:        user.Name,
			Account:     user.Account,
			CreatedTime: user.CreatedTime.String(),
			UpdatedTime: user.UpdatedTime.String(),
			DeletedTime: user.DeletedTime.String(),
		},
	}, nil
}

//ChangePassword 创建.
func (s *SystemAPI) ChangePassword(ctx context.Context, req *v1.ChangePasswordRequest) (*v1.ChangePasswordResponse, error) {
	if req.Password != req.Confirm {
		return nil, status.Error(codes.InvalidArgument, "密码不一致!")
	}
	_, err := system.ChangePassword(ctx, req.Account, common.EncryptPassword(req.Password))
	if err != nil {
		return nil, err
	}
	return &v1.ChangePasswordResponse{
		Code: common.SuccessCode,
		Msg:  common.ResponseMsg(common.SuccessCode),
	}, nil
}
