// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserActionsDao is the data access object for the table user_actions.
type UserActionsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UserActionsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UserActionsColumns defines and stores column names for the table user_actions.
type UserActionsColumns struct {
	Id         string //
	UserId     string //
	Action     string //
	ObjectType string //
	ObjectId   string //
	Extra      string //
	CreatedAt  string //
}

// userActionsColumns holds the columns for the table user_actions.
var userActionsColumns = UserActionsColumns{
	Id:         "id",
	UserId:     "user_id",
	Action:     "action",
	ObjectType: "object_type",
	ObjectId:   "object_id",
	Extra:      "extra",
	CreatedAt:  "created_at",
}

// NewUserActionsDao creates and returns a new DAO object for table data access.
func NewUserActionsDao(handlers ...gdb.ModelHandler) *UserActionsDao {
	return &UserActionsDao{
		group:    "default",
		table:    "user_actions",
		columns:  userActionsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserActionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserActionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserActionsDao) Columns() UserActionsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserActionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserActionsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserActionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
