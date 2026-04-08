package doc

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/doc/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) DocView(ctx context.Context, req *v1.DocViewReq) (res *v1.DocViewRes, err error) {
	if err = service.Doc().IncrementView(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.DocViewRes{}, nil
}
