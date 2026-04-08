package comment

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/comment/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) CommentAdminGetList(ctx context.Context, req *v1.CommentAdminGetListReq) (res *v1.CommentAdminGetListRes, err error) {
	return service.Comment().AdminGetList(ctx, req)
}
