package nav_menu

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/nav_menu/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) NavMenuList(ctx context.Context, req *v1.NavMenuListReq) (*v1.NavMenuListRes, error) {
	list, err := service.NavMenu().List(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.NavMenuListRes{List: list}, nil
}
