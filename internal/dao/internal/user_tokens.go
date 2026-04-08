// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserTokensDao is the data access object for the table user_tokens.
type UserTokensDao struct {
	table    string
	group    string
	columns  UserTokensColumns
	handlers []gdb.ModelHandler
}

// UserTokensColumns defines and stores column names for the table user_tokens.
type UserTokensColumns struct {
	Id         string //
	UserId     string //
	Name       string //
	Prefix     string //
	TokenHash  string //
	ExpiresAt  string //
	LastUsedAt string //
	CreatedAt  string //
}

// userTokensColumns holds the columns for the table user_tokens.
var userTokensColumns = UserTokensColumns{
	Id:         "id",
	UserId:     "user_id",
	Name:       "name",
	Prefix:     "prefix",
	TokenHash:  "token_hash",
	ExpiresAt:  "expires_at",
	LastUsedAt: "last_used_at",
	CreatedAt:  "created_at",
}

// NewUserTokensDao creates and returns a new DAO object for table data access.
func NewUserTokensDao(handlers ...gdb.ModelHandler) *UserTokensDao {
	return &UserTokensDao{
		group:    "default",
		table:    "user_tokens",
		columns:  userTokensColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserTokensDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserTokensDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserTokensDao) Columns() UserTokensColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserTokensDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserTokensDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
func (dao *UserTokensDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
