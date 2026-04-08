package nav_menu

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/nav_menu/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) NavMenuDelete(ctx context.Context, req *v1.NavMenuDeleteReq) (*v1.NavMenuDeleteRes, error) {
	return &v1.NavMenuDeleteRes{}, service.NavMenu().Delete(ctx, req.Id)
}
