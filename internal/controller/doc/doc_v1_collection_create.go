package doc

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/doc/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) CollectionCreate(ctx context.Context, req *v1.CollectionCreateReq) (res *v1.CollectionCreateRes, err error) {
	id, err := service.Doc().CollectionCreate(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.CollectionCreateRes{Id: id}, nil
}
