// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DocsDao is the data access object for the table docs.
type DocsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  DocsColumns        // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// DocsColumns defines and stores column names for the table docs.
type DocsColumns struct {
	Id            string //
	CollectionId  string //
	ParentId      string //
	SortOrder     string //
	Status        string //
	Title         string //
	Slug          string //
	Content       string //
	Excerpt       string //
	AuthorId      string //
	CommentStatus string //
	Locale        string //
	PublishedAt   string //
	CreatedAt     string //
	UpdatedAt     string //
	DeletedAt     string //
}

// docsColumns holds the columns for the table docs.
var docsColumns = DocsColumns{
	Id:            "id",
	CollectionId:  "collection_id",
	ParentId:      "parent_id",
	SortOrder:     "sort_order",
	Status:        "status",
	Title:         "title",
	Slug:          "slug",
	Content:       "content",
	Excerpt:       "excerpt",
	AuthorId:      "author_id",
	CommentStatus: "comment_status",
	Locale:        "locale",
	PublishedAt:   "published_at",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	DeletedAt:     "deleted_at",
}

// NewDocsDao creates and returns a new DAO object for table data access.
func NewDocsDao(handlers ...gdb.ModelHandler) *DocsDao {
	return &DocsDao{
		group:    "default",
		table:    "docs",
		columns:  docsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DocsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DocsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DocsDao) Columns() DocsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DocsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DocsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DocsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
