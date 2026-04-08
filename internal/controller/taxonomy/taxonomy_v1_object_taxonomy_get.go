package taxonomy

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/taxonomy/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) ObjectTaxonomyGet(ctx context.Context, req *v1.ObjectTaxonomyGetReq) (res *v1.ObjectTaxonomyGetRes, err error) {
	return service.Taxonomy().ObjectTaxonomyGet(ctx, req)
}
