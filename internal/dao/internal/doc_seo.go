// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DocSeoDao is the data access object for the table doc_seo.
type DocSeoDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  DocSeoColumns      // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// DocSeoColumns defines and stores column names for the table doc_seo.
type DocSeoColumns struct {
	DocId          string //
	MetaTitle      string //
	MetaDesc       string //
	OgTitle        string //
	OgImage        string //
	CanonicalUrl   string //
	Robots         string //
	StructuredData string //
}

// docSeoColumns holds the columns for the table doc_seo.
var docSeoColumns = DocSeoColumns{
	DocId:          "doc_id",
	MetaTitle:      "meta_title",
	MetaDesc:       "meta_desc",
	OgTitle:        "og_title",
	OgImage:        "og_image",
	CanonicalUrl:   "canonical_url",
	Robots:         "robots",
	StructuredData: "structured_data",
}

// NewDocSeoDao creates and returns a new DAO object for table data access.
func NewDocSeoDao(handlers ...gdb.ModelHandler) *DocSeoDao {
	return &DocSeoDao{
		group:    "default",
		table:    "doc_seo",
		columns:  docSeoColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DocSeoDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DocSeoDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DocSeoDao) Columns() DocSeoColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DocSeoDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DocSeoDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DocSeoDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
