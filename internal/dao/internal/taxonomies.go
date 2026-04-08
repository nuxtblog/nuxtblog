// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TaxonomiesDao is the data access object for the table taxonomies.
type TaxonomiesDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  TaxonomiesColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// TaxonomiesColumns defines and stores column names for the table taxonomies.
type TaxonomiesColumns struct {
	Id          string //
	TermId      string //
	Taxonomy    string //
	Description string //
	ParentId    string //
	PostCount   string //
	Extra       string //
}

// taxonomiesColumns holds the columns for the table taxonomies.
var taxonomiesColumns = TaxonomiesColumns{
	Id:          "id",
	TermId:      "term_id",
	Taxonomy:    "taxonomy",
	Description: "description",
	ParentId:    "parent_id",
	PostCount:   "post_count",
	Extra:       "extra",
}

// NewTaxonomiesDao creates and returns a new DAO object for table data access.
func NewTaxonomiesDao(handlers ...gdb.ModelHandler) *TaxonomiesDao {
	return &TaxonomiesDao{
		group:    "default",
		table:    "taxonomies",
		columns:  taxonomiesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TaxonomiesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TaxonomiesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TaxonomiesDao) Columns() TaxonomiesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TaxonomiesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TaxonomiesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TaxonomiesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
