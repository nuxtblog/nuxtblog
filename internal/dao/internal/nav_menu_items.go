package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type NavMenuItemsDao struct {
	table    string
	group    string
	columns  NavMenuItemsColumns
	handlers []gdb.ModelHandler
}

type NavMenuItemsColumns struct {
	Id         string
	MenuId     string
	ParentId   string
	ObjectType string
	ObjectId   string
	Label      string
	Url        string
	Target     string
	CssClasses string
	SortOrder  string
}

var navMenuItemsColumns = NavMenuItemsColumns{
	Id:         "id",
	MenuId:     "menu_id",
	ParentId:   "parent_id",
	ObjectType: "object_type",
	ObjectId:   "object_id",
	Label:      "label",
	Url:        "url",
	Target:     "target",
	CssClasses: "css_classes",
	SortOrder:  "sort_order",
}

func NewNavMenuItemsDao(handlers ...gdb.ModelHandler) *NavMenuItemsDao {
	return &NavMenuItemsDao{
		group:    "default",
		table:    "nav_menu_items",
		columns:  navMenuItemsColumns,
		handlers: handlers,
	}
}

func (dao *NavMenuItemsDao) DB() gdb.DB                  { return g.DB(dao.group) }
func (dao *NavMenuItemsDao) Table() string                { return dao.table }
func (dao *NavMenuItemsDao) Columns() NavMenuItemsColumns { return dao.columns }
func (dao *NavMenuItemsDao) Group() string                { return dao.group }

func (dao *NavMenuItemsDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, h := range dao.handlers {
		model = h(model)
	}
	return model.Safe().Ctx(ctx)
}

func (dao *NavMenuItemsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) error {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
