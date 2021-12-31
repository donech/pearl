package common

import (
	"context"
	"encoding/base64"
	"errors"
	"flag"
	"strconv"

	"github.com/donech/pearl/internal/data/ent"
	"github.com/donech/pearl/internal/service/system"
	"github.com/donech/pearl/pkg/cipher"
	"github.com/donech/pearl/pkg/jwt"
)

var authKey = "donech"

func init() {
	flag.StringVar(&authKey, "authKey", "donech0123456789", "key for encrypt password")
}

func ValidatePassword(password, encryptPassword string) bool {
	return EncryptPassword(password) == encryptPassword
}

func EncryptPassword(password string) string {
	res := cipher.AesEncryptCBC([]byte(password), []byte(authKey))
	return base64.StdEncoding.EncodeToString(res)
}

//AuthUser 通过解析 jwtToken 获取用户信息
//flag 控制是否去数据库读取详情
func AuthUser(ctx context.Context, flag bool) (*ent.User, error) {
	claims := jwt.GetClaimsFromCtx(ctx)
	if id, ok := claims["id"].(string); ok {
		number, _ := strconv.ParseInt(id, 10, 64)
		if !flag {
			return &ent.User{
				ID:   number,
				Name: claims["name"].(string),
			}, nil
		}
		return system.GetUserById(ctx, number)
	}
	return nil, errors.New("no auth user found")
}
