package follow

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/follow/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *PublicControllerV1) UserFollowing(ctx context.Context, req *v1.UserFollowingReq) (res *v1.UserFollowingRes, err error) {
	return service.Follow().Following(ctx, req)
}
