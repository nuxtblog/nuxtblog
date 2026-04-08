// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MediasDao is the data access object for the table medias.
type MediasDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  MediasColumns      // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// MediasColumns defines and stores column names for the table medias.
type MediasColumns struct {
	Id          string //
	UploaderId  string //
	StorageType string //
	StorageKey  string //
	CdnUrl      string //
	Filename    string //
	MimeType    string //
	FileSize    string //
	Width       string //
	Height      string //
	Duration    string //
	AltText     string //
	Title       string //
	Category    string //
	Variants    string //
	FileMeta    string //
	CreatedAt   string //
	DeletedAt   string //
}

// mediasColumns holds the columns for the table medias.
var mediasColumns = MediasColumns{
	Id:          "id",
	UploaderId:  "uploader_id",
	StorageType: "storage_type",
	StorageKey:  "storage_key",
	CdnUrl:      "cdn_url",
	Filename:    "filename",
	MimeType:    "mime_type",
	FileSize:    "file_size",
	Width:       "width",
	Height:      "height",
	Duration:    "duration",
	AltText:     "alt_text",
	Title:       "title",
	Category:    "category",
	Variants:    "variants",
	FileMeta:    "file_meta",
	CreatedAt:   "created_at",
	DeletedAt:   "deleted_at",
}

// NewMediasDao creates and returns a new DAO object for table data access.
func NewMediasDao(handlers ...gdb.ModelHandler) *MediasDao {
	return &MediasDao{
		group:    "default",
		table:    "medias",
		columns:  mediasColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MediasDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MediasDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MediasDao) Columns() MediasColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MediasDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MediasDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MediasDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
