// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserFollowsDao is the data access object for the table user_follows.
type UserFollowsDao struct {
	table    string
	group    string
	columns  UserFollowsColumns
	handlers []gdb.ModelHandler
}

// UserFollowsColumns defines and stores column names for the table user_follows.
type UserFollowsColumns struct {
	FollowerId  string //
	FollowingId string //
	CreatedAt   string //
}

// userFollowsColumns holds the columns for the table user_follows.
var userFollowsColumns = UserFollowsColumns{
	FollowerId:  "follower_id",
	FollowingId: "following_id",
	CreatedAt:   "created_at",
}

// NewUserFollowsDao creates and returns a new DAO object for table data access.
func NewUserFollowsDao(handlers ...gdb.ModelHandler) *UserFollowsDao {
	return &UserFollowsDao{
		group:    "default",
		table:    "user_follows",
		columns:  userFollowsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserFollowsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserFollowsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserFollowsDao) Columns() UserFollowsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserFollowsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserFollowsDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
func (dao *UserFollowsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
