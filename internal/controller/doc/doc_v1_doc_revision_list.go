package doc

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/doc/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) DocRevisionList(ctx context.Context, req *v1.DocRevisionListReq) (res *v1.DocRevisionListRes, err error) {
	return service.Doc().DocRevisionList(ctx, req)
}
