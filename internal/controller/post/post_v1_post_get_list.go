package post

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/post/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PostGetList(ctx context.Context, req *v1.PostGetListReq) (res *v1.PostGetListRes, err error) {
	return service.Post().GetList(ctx, req)
}
