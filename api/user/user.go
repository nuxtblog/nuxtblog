// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package user

import (
	"context"

	"github.com/nuxtblog/nuxtblog/api/user/v1"
)

type IUserV1 interface {
	UserCreate(ctx context.Context, req *v1.UserCreateReq) (res *v1.UserCreateRes, err error)
	UserDelete(ctx context.Context, req *v1.UserDeleteReq) (res *v1.UserDeleteRes, err error)
	UserUpdate(ctx context.Context, req *v1.UserUpdateReq) (res *v1.UserUpdateRes, err error)
	UserChangePassword(ctx context.Context, req *v1.UserChangePasswordReq) (res *v1.UserChangePasswordRes, err error)
	UserGetOne(ctx context.Context, req *v1.UserGetOneReq) (res *v1.UserGetOneRes, err error)
	UserGetList(ctx context.Context, req *v1.UserGetListReq) (res *v1.UserGetListRes, err error)
}
