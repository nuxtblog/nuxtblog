// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PostsDao is the data access object for the table posts.
type PostsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  PostsColumns       // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// PostsColumns defines and stores column names for the table posts.
type PostsColumns struct {
	Id            string //
	PostType      string //
	Status        string //
	Title         string //
	Slug          string //
	Content       string //
	Excerpt       string //
	AuthorId      string //
	FeaturedImgId string //
	CommentStatus string //
	PasswordHash  string //
	Locale        string //
	PublishedAt   string //
	CreatedAt     string //
	UpdatedAt     string //
	DeletedAt     string //
}

// postsColumns holds the columns for the table posts.
var postsColumns = PostsColumns{
	Id:            "id",
	PostType:      "post_type",
	Status:        "status",
	Title:         "title",
	Slug:          "slug",
	Content:       "content",
	Excerpt:       "excerpt",
	AuthorId:      "author_id",
	FeaturedImgId: "featured_img_id",
	CommentStatus: "comment_status",
	PasswordHash:  "password_hash",
	Locale:        "locale",
	PublishedAt:   "published_at",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	DeletedAt:     "deleted_at",
}

// NewPostsDao creates and returns a new DAO object for table data access.
func NewPostsDao(handlers ...gdb.ModelHandler) *PostsDao {
	return &PostsDao{
		group:    "default",
		table:    "posts",
		columns:  postsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PostsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PostsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PostsDao) Columns() PostsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PostsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PostsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PostsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
