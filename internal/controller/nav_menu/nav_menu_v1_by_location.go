package nav_menu

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/nav_menu/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) NavMenuByLocation(ctx context.Context, req *v1.NavMenuByLocationReq) (*v1.NavMenuByLocationRes, error) {
	out, err := service.NavMenu().GetByLocation(ctx, req.Location)
	if err != nil {
		return nil, err
	}
	return &v1.NavMenuByLocationRes{NavMenuOutput: out}, nil
}
