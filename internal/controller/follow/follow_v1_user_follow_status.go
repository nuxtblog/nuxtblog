package follow

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/follow/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *AuthControllerV1) UserFollowStatus(ctx context.Context, req *v1.UserFollowStatusReq) (res *v1.UserFollowStatusRes, err error) {
	return service.Follow().FollowStatus(ctx, req)
}
