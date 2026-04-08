package follow

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/follow/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *AuthControllerV1) UserRemoveFollower(ctx context.Context, req *v1.UserRemoveFollowerReq) (res *v1.UserRemoveFollowerRes, err error) {
	return service.Follow().RemoveFollower(ctx, req)
}
