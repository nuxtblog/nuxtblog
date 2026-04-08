// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PostStatsDao is the data access object for the table post_stats.
type PostStatsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  PostStatsColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// PostStatsColumns defines and stores column names for the table post_stats.
type PostStatsColumns struct {
	PostId       string //
	ViewCount    string //
	LikeCount    string //
	CommentCount string //
	ShareCount   string //
	UpdatedAt    string //
}

// postStatsColumns holds the columns for the table post_stats.
var postStatsColumns = PostStatsColumns{
	PostId:       "post_id",
	ViewCount:    "view_count",
	LikeCount:    "like_count",
	CommentCount: "comment_count",
	ShareCount:   "share_count",
	UpdatedAt:    "updated_at",
}

// NewPostStatsDao creates and returns a new DAO object for table data access.
func NewPostStatsDao(handlers ...gdb.ModelHandler) *PostStatsDao {
	return &PostStatsDao{
		group:    "default",
		table:    "post_stats",
		columns:  postStatsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PostStatsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PostStatsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PostStatsDao) Columns() PostStatsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PostStatsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PostStatsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PostStatsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
