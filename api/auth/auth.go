// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package auth

import (
	"context"

	"github.com/nuxtblog/nuxtblog/api/auth/v1"
)

type IAuthV1 interface {
	AuthLogin(ctx context.Context, req *v1.AuthLoginReq) (res *v1.AuthLoginRes, err error)
	AuthRegister(ctx context.Context, req *v1.AuthRegisterReq) (res *v1.AuthRegisterRes, err error)
	AuthLogout(ctx context.Context, req *v1.AuthLogoutReq) (res *v1.AuthLogoutRes, err error)
	AuthRefresh(ctx context.Context, req *v1.AuthRefreshReq) (res *v1.AuthRefreshRes, err error)
	AuthMe(ctx context.Context, req *v1.AuthMeReq) (res *v1.AuthMeRes, err error)
	AuthOAuthProviders(ctx context.Context, req *v1.AuthOAuthProvidersReq) (res *v1.AuthOAuthProvidersRes, err error)
}
