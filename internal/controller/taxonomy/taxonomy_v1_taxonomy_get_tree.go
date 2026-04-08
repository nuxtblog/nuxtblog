package taxonomy

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/taxonomy/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) TaxonomyGetTree(ctx context.Context, req *v1.TaxonomyGetTreeReq) (res *v1.TaxonomyGetTreeRes, err error) {
	return service.Taxonomy().TaxonomyGetTree(ctx, req.Taxonomy)
}
