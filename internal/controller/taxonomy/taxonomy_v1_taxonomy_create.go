package taxonomy

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/taxonomy/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) TaxonomyCreate(ctx context.Context, req *v1.TaxonomyCreateReq) (res *v1.TaxonomyCreateRes, err error) {
	id, err := service.Taxonomy().TaxonomyCreate(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.TaxonomyCreateRes{Id: id}, nil
}
