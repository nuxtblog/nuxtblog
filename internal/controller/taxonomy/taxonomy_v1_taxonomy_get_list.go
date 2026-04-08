package taxonomy

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/taxonomy/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) TaxonomyGetList(ctx context.Context, req *v1.TaxonomyGetListReq) (res *v1.TaxonomyGetListRes, err error) {
	return service.Taxonomy().TaxonomyGetList(ctx, req)
}
