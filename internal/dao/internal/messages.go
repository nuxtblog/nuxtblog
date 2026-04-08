// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MessagesDao is the data access object for the table messages.
type MessagesDao struct {
	table    string
	group    string
	columns  MessagesColumns
	handlers []gdb.ModelHandler
}

// MessagesColumns defines and stores column names for the table messages.
type MessagesColumns struct {
	Id             string //
	ConversationId string //
	SenderId       string //
	Content        string //
	IsRead         string //
	CreatedAt      string //
}

// messagesColumns holds the columns for the table messages.
var messagesColumns = MessagesColumns{
	Id:             "id",
	ConversationId: "conversation_id",
	SenderId:       "sender_id",
	Content:        "content",
	IsRead:         "is_read",
	CreatedAt:      "created_at",
}

// NewMessagesDao creates and returns a new DAO object for table data access.
func NewMessagesDao(handlers ...gdb.ModelHandler) *MessagesDao {
	return &MessagesDao{
		group:    "default",
		table:    "messages",
		columns:  messagesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MessagesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MessagesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MessagesDao) Columns() MessagesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MessagesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MessagesDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
func (dao *MessagesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
