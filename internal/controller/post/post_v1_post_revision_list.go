package post

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/post/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PostRevisionList(ctx context.Context, req *v1.PostRevisionListReq) (res *v1.PostRevisionListRes, err error) {
	return service.Post().RevisionList(ctx, req)
}
