package service

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/nav_menu/v1"
)

type INavMenu interface {
	List(ctx context.Context) ([]*v1.NavMenuOutput, error)
	Create(ctx context.Context, req *v1.NavMenuCreateReq) (*v1.NavMenuOutput, error)
	GetOne(ctx context.Context, id int64) (*v1.NavMenuOutput, error)
	Update(ctx context.Context, req *v1.NavMenuUpdateReq) (*v1.NavMenuOutput, error)
	Delete(ctx context.Context, id int64) error
	GetByLocation(ctx context.Context, location string) (*v1.NavMenuOutput, error)
}

var localNavMenu INavMenu

func NavMenu() INavMenu {
	if localNavMenu == nil {
		panic("implement not found for interface INavMenu, forgot register?")
	}
	return localNavMenu
}

func RegisterNavMenu(i INavMenu) { localNavMenu = i }
