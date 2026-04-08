package taxonomy

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/taxonomy/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) TaxonomyDelete(ctx context.Context, req *v1.TaxonomyDeleteReq) (res *v1.TaxonomyDeleteRes, err error) {
	return &v1.TaxonomyDeleteRes{}, service.Taxonomy().TaxonomyDelete(ctx, req.Id)
}
