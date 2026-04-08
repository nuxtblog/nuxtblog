package service

import (
	v1 "github.com/nuxtblog/nuxtblog/api/auth/v1"
	"context"
)

type (
	IAuth interface {
		Login(ctx context.Context, req *v1.AuthLoginReq) (*v1.AuthLoginRes, error)
		Register(ctx context.Context, req *v1.AuthRegisterReq) (*v1.AuthRegisterRes, error)
		Logout(ctx context.Context, refreshToken string) error
		Me(ctx context.Context) (*v1.AuthUserItem, error)
		Refresh(ctx context.Context, refreshToken string) (*v1.AuthRefreshRes, error)
		// OAuth
		OAuthProviders(ctx context.Context) []string
		OAuthRedirect(ctx context.Context, provider, redirectAfter string) (string, error)
		OAuthCallback(ctx context.Context, provider, code, state string) (access, refresh string, expiresIn int64, redirectAfter string, err error)
	}
)

var localAuth IAuth

func Auth() IAuth {
	if localAuth == nil {
		panic("implement not found for interface IAuth, forgot register?")
	}
	return localAuth
}

func RegisterAuth(i IAuth) {
	localAuth = i
}
