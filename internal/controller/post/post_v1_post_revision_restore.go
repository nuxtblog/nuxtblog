package post

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/post/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PostRevisionRestore(ctx context.Context, req *v1.PostRevisionRestoreReq) (res *v1.PostRevisionRestoreRes, err error) {
	return &v1.PostRevisionRestoreRes{}, service.Post().RevisionRestore(ctx, req.Id, req.RevisionId)
}
