package dao

import "github.com/nuxtblog/nuxtblog/internal/dao/internal"

type navMenusDao struct {
	*internal.NavMenusDao
}

var NavMenus = navMenusDao{internal.NewNavMenusDao()}
