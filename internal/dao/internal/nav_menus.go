package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type NavMenusDao struct {
	table    string
	group    string
	columns  NavMenusColumns
	handlers []gdb.ModelHandler
}

type NavMenusColumns struct {
	Id          string
	Name        string
	Location    string
	Description string
	CreatedAt   string
	UpdatedAt   string
}

var navMenusColumns = NavMenusColumns{
	Id:          "id",
	Name:        "name",
	Location:    "location",
	Description: "description",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

func NewNavMenusDao(handlers ...gdb.ModelHandler) *NavMenusDao {
	return &NavMenusDao{
		group:    "default",
		table:    "nav_menus",
		columns:  navMenusColumns,
		handlers: handlers,
	}
}

func (dao *NavMenusDao) DB() gdb.DB               { return g.DB(dao.group) }
func (dao *NavMenusDao) Table() string             { return dao.table }
func (dao *NavMenusDao) Columns() NavMenusColumns  { return dao.columns }
func (dao *NavMenusDao) Group() string             { return dao.group }

func (dao *NavMenusDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, h := range dao.handlers {
		model = h(model)
	}
	return model.Safe().Ctx(ctx)
}

func (dao *NavMenusDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) error {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
