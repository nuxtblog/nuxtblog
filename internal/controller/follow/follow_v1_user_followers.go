package follow

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/follow/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *PublicControllerV1) UserFollowers(ctx context.Context, req *v1.UserFollowersReq) (res *v1.UserFollowersRes, err error) {
	return service.Follow().Followers(ctx, req)
}
