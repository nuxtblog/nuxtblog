package doc

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/doc/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) CollectionGetOne(ctx context.Context, req *v1.CollectionGetOneReq) (res *v1.CollectionGetOneRes, err error) {
	item, err := service.Doc().CollectionGetOne(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "collection not found")
	}
	return &v1.CollectionGetOneRes{DocCollectionItem: item}, nil
}
