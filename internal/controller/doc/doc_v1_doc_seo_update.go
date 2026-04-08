package doc

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/doc/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) DocSeoUpdate(ctx context.Context, req *v1.DocSeoUpdateReq) (res *v1.DocSeoUpdateRes, err error) {
	if err = service.Doc().DocSeoUpdate(ctx, req); err != nil {
		return nil, err
	}
	return &v1.DocSeoUpdateRes{}, nil
}
