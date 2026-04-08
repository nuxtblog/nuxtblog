package follow

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/follow/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *AuthControllerV1) UserFollow(ctx context.Context, req *v1.UserFollowReq) (res *v1.UserFollowRes, err error) {
	return service.Follow().Follow(ctx, req)
}
