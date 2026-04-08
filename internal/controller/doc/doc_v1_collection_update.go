package doc

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/doc/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) CollectionUpdate(ctx context.Context, req *v1.CollectionUpdateReq) (res *v1.CollectionUpdateRes, err error) {
	if err = service.Doc().CollectionUpdate(ctx, req); err != nil {
		return nil, err
	}
	return &v1.CollectionUpdateRes{}, nil
}
