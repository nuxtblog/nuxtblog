package follow

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/follow/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *AuthControllerV1) UserUnfollow(ctx context.Context, req *v1.UserUnfollowReq) (res *v1.UserUnfollowRes, err error) {
	return service.Follow().Unfollow(ctx, req)
}
