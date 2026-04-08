package dao

import (
	"github.com/nuxtblog/nuxtblog/internal/dao/internal"
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type notificationsDao struct {
	*internal.NotificationsDao
}

var Notifications = notificationsDao{internal.NewNotificationsDao()}

func (d *notificationsDao) Create(ctx context.Context, data g.Map) (int64, error) {
	result, err := d.Ctx(ctx).Insert(data)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (d *notificationsDao) UnreadCount(ctx context.Context, userId int64) (int, error) {
	return d.Ctx(ctx).
		Where("user_id", userId).
		Where("is_read", 0).
		WhereNull("deleted_at").
		Count()
}

func (d *notificationsDao) MarkRead(ctx context.Context, id int64) error {
	_, err := d.Ctx(ctx).Where("id", id).WhereNull("deleted_at").Update(g.Map{"is_read": 1})
	return err
}

func (d *notificationsDao) MarkAllRead(ctx context.Context, userId int64) error {
	_, err := d.Ctx(ctx).
		Where("user_id", userId).
		Where("is_read", 0).
		WhereNull("deleted_at").
		Update(g.Map{"is_read": 1})
	return err
}

func (d *notificationsDao) SoftDelete(ctx context.Context, id int64) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := d.Ctx(ctx).Where("id", id).WhereNull("deleted_at").Update(g.Map{"deleted_at": now})
	return err
}

// ClearByFilter soft-deletes notifications by filter: "all", "unread", "interaction", "system"
func (d *notificationsDao) ClearByFilter(ctx context.Context, userId int64, filter string) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	m := d.Ctx(ctx).Where("user_id", userId).WhereNull("deleted_at")
	switch filter {
	case "unread":
		m = m.Where("is_read", 0)
	case "interaction":
		m = m.WhereNotIn("type", g.Slice{"system"})
	case "system":
		m = m.Where("type", "system")
	}
	_, err := m.Update(g.Map{"deleted_at": now})
	return err
}
