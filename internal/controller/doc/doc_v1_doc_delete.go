package doc

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/doc/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) DocDelete(ctx context.Context, req *v1.DocDeleteReq) (res *v1.DocDeleteRes, err error) {
	if err = service.Doc().DocDelete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.DocDeleteRes{}, nil
}
