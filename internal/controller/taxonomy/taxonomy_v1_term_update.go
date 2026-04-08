package taxonomy

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/taxonomy/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) TermUpdate(ctx context.Context, req *v1.TermUpdateReq) (res *v1.TermUpdateRes, err error) {
	return &v1.TermUpdateRes{}, service.Taxonomy().TermUpdate(ctx, req)
}
