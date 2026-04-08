package taxonomy

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/taxonomy/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) TaxonomyUpdate(ctx context.Context, req *v1.TaxonomyUpdateReq) (res *v1.TaxonomyUpdateRes, err error) {
	return &v1.TaxonomyUpdateRes{}, service.Taxonomy().TaxonomyUpdate(ctx, req)
}
