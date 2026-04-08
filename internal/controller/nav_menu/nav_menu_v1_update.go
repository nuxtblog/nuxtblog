package nav_menu

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/nav_menu/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) NavMenuUpdate(ctx context.Context, req *v1.NavMenuUpdateReq) (*v1.NavMenuUpdateRes, error) {
	out, err := service.NavMenu().Update(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.NavMenuUpdateRes{NavMenuOutput: out}, nil
}
