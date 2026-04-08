// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserLikesDao is the data access object for the table user_likes.
type UserLikesDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UserLikesColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UserLikesColumns defines and stores column names for the table user_likes.
type UserLikesColumns struct {
	UserId     string //
	ObjectType string //
	ObjectId   string //
	CreatedAt  string //
}

// userLikesColumns holds the columns for the table user_likes.
var userLikesColumns = UserLikesColumns{
	UserId:     "user_id",
	ObjectType: "object_type",
	ObjectId:   "object_id",
	CreatedAt:  "created_at",
}

// NewUserLikesDao creates and returns a new DAO object for table data access.
func NewUserLikesDao(handlers ...gdb.ModelHandler) *UserLikesDao {
	return &UserLikesDao{
		group:    "default",
		table:    "user_likes",
		columns:  userLikesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserLikesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserLikesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserLikesDao) Columns() UserLikesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserLikesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserLikesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserLikesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
