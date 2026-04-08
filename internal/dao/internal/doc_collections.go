// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DocCollectionsDao is the data access object for the table doc_collections.
type DocCollectionsDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  DocCollectionsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// DocCollectionsColumns defines and stores column names for the table doc_collections.
type DocCollectionsColumns struct {
	Id          string //
	Slug        string //
	Title       string //
	Description string //
	CoverImgId  string //
	AuthorId    string //
	Status      string //
	Locale      string //
	SortOrder   string //
	CreatedAt   string //
	UpdatedAt   string //
	DeletedAt   string //
}

// docCollectionsColumns holds the columns for the table doc_collections.
var docCollectionsColumns = DocCollectionsColumns{
	Id:          "id",
	Slug:        "slug",
	Title:       "title",
	Description: "description",
	CoverImgId:  "cover_img_id",
	AuthorId:    "author_id",
	Status:      "status",
	Locale:      "locale",
	SortOrder:   "sort_order",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewDocCollectionsDao creates and returns a new DAO object for table data access.
func NewDocCollectionsDao(handlers ...gdb.ModelHandler) *DocCollectionsDao {
	return &DocCollectionsDao{
		group:    "default",
		table:    "doc_collections",
		columns:  docCollectionsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DocCollectionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DocCollectionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DocCollectionsDao) Columns() DocCollectionsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DocCollectionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DocCollectionsDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *DocCollectionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
