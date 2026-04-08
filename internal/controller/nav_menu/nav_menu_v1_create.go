package nav_menu

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/nav_menu/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) NavMenuCreate(ctx context.Context, req *v1.NavMenuCreateReq) (*v1.NavMenuCreateRes, error) {
	out, err := service.NavMenu().Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.NavMenuCreateRes{NavMenuOutput: out}, nil
}
