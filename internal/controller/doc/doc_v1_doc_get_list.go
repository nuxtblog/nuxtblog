package doc

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/doc/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) DocGetList(ctx context.Context, req *v1.DocGetListReq) (res *v1.DocGetListRes, err error) {
	return service.Doc().DocGetList(ctx, req)
}
