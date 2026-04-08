package doc

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/doc/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) DocUpdate(ctx context.Context, req *v1.DocUpdateReq) (res *v1.DocUpdateRes, err error) {
	if err = service.Doc().DocUpdate(ctx, req); err != nil {
		return nil, err
	}
	return &v1.DocUpdateRes{}, nil
}
