package taxonomy

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/taxonomy/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) ObjectTaxonomyUnbind(ctx context.Context, req *v1.ObjectTaxonomyUnbindReq) (res *v1.ObjectTaxonomyUnbindRes, err error) {
	return &v1.ObjectTaxonomyUnbindRes{}, service.Taxonomy().ObjectTaxonomyUnbind(ctx, req)
}
