package service

import (
	v1 "github.com/nuxtblog/nuxtblog/api/user/v1"
	"context"
)

type (
	IUser interface {
		GetList(ctx context.Context, req *v1.UserGetListReq) (*v1.UserGetListRes, error)
		GetOne(ctx context.Context, id int64) (*v1.UserItem, error)
		Create(ctx context.Context, req *v1.UserCreateReq) (int64, error)
		Update(ctx context.Context, req *v1.UserUpdateReq) error
		Delete(ctx context.Context, id int64) error
		ChangePassword(ctx context.Context, req *v1.UserChangePasswordReq) error
	}
)

var (
	localUser IUser
)

func User() IUser {
	if localUser == nil {
		panic("implement not found for interface IUser, forgot register?")
	}
	return localUser
}

func RegisterUser(i IUser) {
	localUser = i
}
