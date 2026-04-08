package doc

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/doc/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) CollectionGetList(ctx context.Context, req *v1.CollectionGetListReq) (res *v1.CollectionGetListRes, err error) {
	return service.Doc().CollectionGetList(ctx, req)
}
