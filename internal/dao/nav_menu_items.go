package dao

import "github.com/nuxtblog/nuxtblog/internal/dao/internal"

type navMenuItemsDao struct {
	*internal.NavMenuItemsDao
}

var NavMenuItems = navMenuItemsDao{internal.NewNavMenuItemsDao()}
