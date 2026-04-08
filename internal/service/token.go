package service

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/token/v1"
)

type IToken interface {
	List(ctx context.Context, req *v1.UserTokenListReq) (*v1.UserTokenListRes, error)
	Create(ctx context.Context, req *v1.UserTokenCreateReq) (*v1.UserTokenCreateRes, error)
	Delete(ctx context.Context, req *v1.UserTokenDeleteReq) (*v1.UserTokenDeleteRes, error)
}

var localToken IToken

func Token() IToken {
	if localToken == nil {
		panic("implement not found for interface IToken, forgot register?")
	}
	return localToken
}

func RegisterToken(i IToken) {
	localToken = i
}
