package system

import (
	"context"

	"github.com/donech/pearl/internal/data"
	"github.com/donech/pearl/internal/data/ent"
	"github.com/donech/pearl/internal/data/ent/predicate"
	"github.com/donech/pearl/internal/data/ent/user"
)

type UsersParams struct {
	Name    string
	Account string
	Page    int
	Size    int
}

type UserList struct {
	Users []*ent.User
	Total int
}

//Users 获取 system user 列表
//condition[0] query  string
//condition[1:] the args for query string
func Users(ctx context.Context, params UsersParams) (*UserList, error) {
	if params.Page <= 0 {
		params.Page = 1
	}

	if params.Size <= 0 {
		params.Size = 10
	}

	where := make([]predicate.User, 0)
	if len(params.Name) > 0 {
		where = append(where, user.NameContains(params.Name))
	}
	if len(params.Account) > 0 {
		where = append(where, user.AccountContains(params.Account))
	}

	builder := data.DB.User.Query().Unique(false).Where(where...)
	res := new(UserList)
	total, err := builder.Clone().Count(ctx)
	if err != nil {
		return res, err
	}
	res.Total = total
	users, err := builder.Offset((params.Page - 1) * params.Size).Limit(params.Size).Order(ent.Desc(user.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	res.Users = users
	return res, nil
}

func SaveUser(ctx context.Context, account, name, password string) (*ent.User, error) {
	userM, err := GetUserByAccount(ctx, account)
	if ent.IsNotFound(err) {
		userM, err := data.DB.User.Create().SetName(name).SetAccount(account).SetPassword(password).Save(ctx)
		if err != nil {
			return nil, err
		}
		return userM, nil
	}
	if err != nil {
		return nil, err
	}
	_, err = data.DB.User.Update().SetName(name).SetPassword(password).Where(user.IDEQ(userM.ID)).Save(ctx)
	return userM, err
}

func ChangePassword(ctx context.Context, account, password string) (*ent.User, error) {
	userM, err := GetUserByAccount(ctx, account)
	if err != nil {
		return nil, err
	}
	_, err = data.DB.User.Update().SetPassword(password).Where(user.IDEQ(userM.ID)).Save(ctx)
	return userM, err
}

func GetUserByAccount(ctx context.Context, account string) (*ent.User, error) {
	return data.DB.User.Query().Where(user.Account(account)).First(ctx)
}

func GetUserById(ctx context.Context, id int64) (*ent.User, error) {
	return data.DB.User.Query().Where(user.ID(id)).First(ctx)
}
