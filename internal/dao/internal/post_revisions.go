// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PostRevisionsDao is the data access object for the table post_revisions.
type PostRevisionsDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  PostRevisionsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// PostRevisionsColumns defines and stores column names for the table post_revisions.
type PostRevisionsColumns struct {
	Id        string //
	PostId    string //
	AuthorId  string //
	Title     string //
	Content   string //
	RevNote   string //
	CreatedAt string //
}

// postRevisionsColumns holds the columns for the table post_revisions.
var postRevisionsColumns = PostRevisionsColumns{
	Id:        "id",
	PostId:    "post_id",
	AuthorId:  "author_id",
	Title:     "title",
	Content:   "content",
	RevNote:   "rev_note",
	CreatedAt: "created_at",
}

// NewPostRevisionsDao creates and returns a new DAO object for table data access.
func NewPostRevisionsDao(handlers ...gdb.ModelHandler) *PostRevisionsDao {
	return &PostRevisionsDao{
		group:    "default",
		table:    "post_revisions",
		columns:  postRevisionsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PostRevisionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PostRevisionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PostRevisionsDao) Columns() PostRevisionsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PostRevisionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PostRevisionsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PostRevisionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
