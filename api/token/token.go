package token

import (
	"context"

	"github.com/nuxtblog/nuxtblog/api/token/v1"
)

type ITokenV1 interface {
	UserTokenList(ctx context.Context, req *v1.UserTokenListReq) (res *v1.UserTokenListRes, err error)
	UserTokenCreate(ctx context.Context, req *v1.UserTokenCreateReq) (res *v1.UserTokenCreateRes, err error)
	UserTokenDelete(ctx context.Context, req *v1.UserTokenDeleteReq) (res *v1.UserTokenDeleteRes, err error)
}
