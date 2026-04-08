package auth

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/auth/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) AuthOAuthProviders(ctx context.Context, req *v1.AuthOAuthProvidersReq) (res *v1.AuthOAuthProvidersRes, err error) {
	providers := service.Auth().OAuthProviders(ctx)
	if providers == nil {
		providers = []string{}
	}
	return &v1.AuthOAuthProvidersRes{Providers: providers}, nil
}
