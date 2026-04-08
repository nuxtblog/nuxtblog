package doc

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/doc/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) DocCreate(ctx context.Context, req *v1.DocCreateReq) (res *v1.DocCreateRes, err error) {
	id, err := service.Doc().DocCreate(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.DocCreateRes{Id: id}, nil
}
