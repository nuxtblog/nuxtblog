package taxonomy

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/taxonomy/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) TermCreate(ctx context.Context, req *v1.TermCreateReq) (res *v1.TermCreateRes, err error) {
	id, err := service.Taxonomy().TermCreate(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.TermCreateRes{Id: id}, nil
}
