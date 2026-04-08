// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// VerificationCodesDao is the data access object for the table verification_codes.
type VerificationCodesDao struct {
	table    string
	group    string
	columns  VerificationCodesColumns
	handlers []gdb.ModelHandler
}

// VerificationCodesColumns defines and stores column names for the table verification_codes.
type VerificationCodesColumns struct {
	Id        string //
	Target    string //
	Code      string //
	Type      string //
	ExpiresAt string //
	UsedAt    string //
	CreatedAt string //
}

// verificationCodesColumns holds the columns for the table verification_codes.
var verificationCodesColumns = VerificationCodesColumns{
	Id:        "id",
	Target:    "target",
	Code:      "code",
	Type:      "type",
	ExpiresAt: "expires_at",
	UsedAt:    "used_at",
	CreatedAt: "created_at",
}

// NewVerificationCodesDao creates and returns a new DAO object for table data access.
func NewVerificationCodesDao(handlers ...gdb.ModelHandler) *VerificationCodesDao {
	return &VerificationCodesDao{
		group:    "default",
		table:    "verification_codes",
		columns:  verificationCodesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *VerificationCodesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *VerificationCodesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *VerificationCodesDao) Columns() VerificationCodesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *VerificationCodesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *VerificationCodesDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
func (dao *VerificationCodesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
