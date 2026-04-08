// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MomentStatsDao is the data access object for the table moment_stats.
type MomentStatsDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  MomentStatsColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// MomentStatsColumns defines and stores column names for the table moment_stats.
type MomentStatsColumns struct {
	MomentId     string //
	ViewCount    string //
	LikeCount    string //
	CommentCount string //
	UpdatedAt    string //
}

// momentStatsColumns holds the columns for the table moment_stats.
var momentStatsColumns = MomentStatsColumns{
	MomentId:     "moment_id",
	ViewCount:    "view_count",
	LikeCount:    "like_count",
	CommentCount: "comment_count",
	UpdatedAt:    "updated_at",
}

// NewMomentStatsDao creates and returns a new DAO object for table data access.
func NewMomentStatsDao(handlers ...gdb.ModelHandler) *MomentStatsDao {
	return &MomentStatsDao{
		group:    "default",
		table:    "moment_stats",
		columns:  momentStatsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MomentStatsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MomentStatsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MomentStatsDao) Columns() MomentStatsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MomentStatsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MomentStatsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MomentStatsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
