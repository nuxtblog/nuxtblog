package taxonomy

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/taxonomy/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) ObjectTaxonomyBind(ctx context.Context, req *v1.ObjectTaxonomyBindReq) (res *v1.ObjectTaxonomyBindRes, err error) {
	return &v1.ObjectTaxonomyBindRes{}, service.Taxonomy().ObjectTaxonomyBind(ctx, req)
}
