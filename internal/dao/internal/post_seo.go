// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PostSeoDao is the data access object for the table post_seo.
type PostSeoDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  PostSeoColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// PostSeoColumns defines and stores column names for the table post_seo.
type PostSeoColumns struct {
	PostId         string //
	MetaTitle      string //
	MetaDesc       string //
	OgTitle        string //
	OgImage        string //
	CanonicalUrl   string //
	Robots         string //
	StructuredData string //
}

// postSeoColumns holds the columns for the table post_seo.
var postSeoColumns = PostSeoColumns{
	PostId:         "post_id",
	MetaTitle:      "meta_title",
	MetaDesc:       "meta_desc",
	OgTitle:        "og_title",
	OgImage:        "og_image",
	CanonicalUrl:   "canonical_url",
	Robots:         "robots",
	StructuredData: "structured_data",
}

// NewPostSeoDao creates and returns a new DAO object for table data access.
func NewPostSeoDao(handlers ...gdb.ModelHandler) *PostSeoDao {
	return &PostSeoDao{
		group:    "default",
		table:    "post_seo",
		columns:  postSeoColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PostSeoDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PostSeoDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PostSeoDao) Columns() PostSeoColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PostSeoDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PostSeoDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PostSeoDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
