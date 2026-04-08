package doc

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/doc/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) CollectionDelete(ctx context.Context, req *v1.CollectionDeleteReq) (res *v1.CollectionDeleteRes, err error) {
	if err = service.Doc().CollectionDelete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.CollectionDeleteRes{}, nil
}
