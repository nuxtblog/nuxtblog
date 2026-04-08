package taxonomy

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/taxonomy/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) TermDelete(ctx context.Context, req *v1.TermDeleteReq) (res *v1.TermDeleteRes, err error) {
	return &v1.TermDeleteRes{}, service.Taxonomy().TermDelete(ctx, req.Id)
}
