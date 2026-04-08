package nav_menu

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/nav_menu/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) NavMenuGetOne(ctx context.Context, req *v1.NavMenuGetOneReq) (*v1.NavMenuGetOneRes, error) {
	out, err := service.NavMenu().GetOne(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.NavMenuGetOneRes{NavMenuOutput: out}, nil
}
