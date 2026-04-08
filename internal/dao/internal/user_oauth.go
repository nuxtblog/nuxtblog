// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserOauthDao is the data access object for the table user_oauth.
type UserOauthDao struct {
	table    string
	group    string
	columns  UserOauthColumns
	handlers []gdb.ModelHandler
}

// UserOauthColumns defines and stores column names for the table user_oauth.
type UserOauthColumns struct {
	Id         string //
	UserId     string //
	Provider   string //
	ProviderId string //
	CreatedAt  string //
}

// userOauthColumns holds the columns for the table user_oauth.
var userOauthColumns = UserOauthColumns{
	Id:         "id",
	UserId:     "user_id",
	Provider:   "provider",
	ProviderId: "provider_id",
	CreatedAt:  "created_at",
}

// NewUserOauthDao creates and returns a new DAO object for table data access.
func NewUserOauthDao(handlers ...gdb.ModelHandler) *UserOauthDao {
	return &UserOauthDao{
		group:    "default",
		table:    "user_oauth",
		columns:  userOauthColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserOauthDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserOauthDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserOauthDao) Columns() UserOauthColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserOauthDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserOauthDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
func (dao *UserOauthDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
