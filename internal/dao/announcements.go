package dao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type announcementsDao struct{}

var Announcements = &announcementsDao{}

// AnnouncementRow is the typed scan target for announcement rows.
type AnnouncementRow struct {
	Id        int64  `orm:"id"`
	Title     string `orm:"title"`
	Content   string `orm:"content"`
	Type      string `orm:"type"`
	CreatedAt string `orm:"created_at"`
}

func (d *announcementsDao) Ctx(ctx context.Context) *gdb.Model {
	return g.DB().Model("announcements").Ctx(ctx)
}

func (d *announcementsDao) Create(ctx context.Context, data g.Map) (int64, error) {
	result, err := d.Ctx(ctx).Insert(data)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (d *announcementsDao) Update(ctx context.Context, id int64, data g.Map) error {
	_, err := d.Ctx(ctx).Where("id", id).WhereNull("deleted_at").Update(data)
	return err
}

func (d *announcementsDao) SoftDelete(ctx context.Context, id int64) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := d.Ctx(ctx).Where("id", id).WhereNull("deleted_at").Update(g.Map{"deleted_at": now})
	return err
}

func (d *announcementsDao) List(ctx context.Context, page, size int) ([]AnnouncementRow, int, error) {
	m := d.Ctx(ctx).WhereNull("deleted_at").OrderDesc("created_at")
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var rows []AnnouncementRow
	if total > 0 {
		err = m.Page(page, size).Scan(&rows)
	}
	return rows, total, err
}

// UnreadCount returns the number of announcements created after the user's cursor.
// If the user has no cursor row, all announcements are considered unread.
func (d *announcementsDao) UnreadCount(ctx context.Context, userId int64) (int, error) {
	var cursor *struct {
		LastReadAt string `orm:"last_read_at"`
	}
	_ = g.DB().Model("user_announcement_cursor").Ctx(ctx).
		Where("user_id", userId).
		Fields("last_read_at").
		Scan(&cursor)

	m := d.Ctx(ctx).WhereNull("deleted_at")
	if cursor != nil && cursor.LastReadAt != "" {
		m = m.WhereGT("created_at", cursor.LastReadAt)
	}
	return m.Count()
}

// UpsertCursor sets or updates the user's last_read_at timestamp.
func (d *announcementsDao) UpsertCursor(ctx context.Context, userId int64, t time.Time) error {
	ts := t.Format("2006-01-02 15:04:05.000")
	_, err := g.DB().Exec(ctx,
		"INSERT INTO user_announcement_cursor(user_id, last_read_at) VALUES(?,?) "+
			"ON CONFLICT(user_id) DO UPDATE SET last_read_at=excluded.last_read_at",
		userId, ts,
	)
	return err
}

// GetCursor returns the user's last_read_at, or nil if no cursor row exists.
func (d *announcementsDao) GetCursor(ctx context.Context, userId int64) (*time.Time, error) {
	var row *struct {
		LastReadAt string `orm:"last_read_at"`
	}
	err := g.DB().Model("user_announcement_cursor").Ctx(ctx).
		Where("user_id", userId).
		Fields("last_read_at").
		Scan(&row)
	if err != nil || row == nil || row.LastReadAt == "" {
		return nil, err
	}
	t, err := time.Parse("2006-01-02 15:04:05", row.LastReadAt)
	if err != nil {
		t, err = time.Parse("2006-01-02 15:04:05.000", row.LastReadAt)
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}
