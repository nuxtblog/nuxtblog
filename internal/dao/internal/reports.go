// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ReportsDao is the data access object for the table reports.
type ReportsDao struct {
	table    string
	group    string
	columns  ReportsColumns
	handlers []gdb.ModelHandler
}

// ReportsColumns defines and stores column names for the table reports.
type ReportsColumns struct {
	Id         string //
	ReporterId string //
	TargetType string //
	TargetId   string //
	Reason     string //
	Detail     string //
	Status     string //
	Notes      string //
	CreatedAt  string //
	ResolvedAt string //
}

// reportsColumns holds the columns for the table reports.
var reportsColumns = ReportsColumns{
	Id:         "id",
	ReporterId: "reporter_id",
	TargetType: "target_type",
	TargetId:   "target_id",
	Reason:     "reason",
	Detail:     "detail",
	Status:     "status",
	Notes:      "notes",
	CreatedAt:  "created_at",
	ResolvedAt: "resolved_at",
}

// NewReportsDao creates and returns a new DAO object for table data access.
func NewReportsDao(handlers ...gdb.ModelHandler) *ReportsDao {
	return &ReportsDao{
		group:    "default",
		table:    "reports",
		columns:  reportsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ReportsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ReportsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ReportsDao) Columns() ReportsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ReportsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ReportsDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
func (dao *ReportsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
