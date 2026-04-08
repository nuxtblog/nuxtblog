package announcement

import (
	"context"
	"time"

	v1 "github.com/nuxtblog/nuxtblog/api/announcement/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/nuxtblog/nuxtblog/internal/util/idgen"
	"github.com/gogf/gf/v2/frame/g"
)

type sAnnouncement struct{}

func New() service.IAnnouncement { return &sAnnouncement{} }

func init() {
	service.RegisterAnnouncement(New())
}

func (s *sAnnouncement) Create(ctx context.Context, title, content, atype string, createdBy int64) (int64, error) {
	if len(content) > 5000 {
		content = content[:5000]
	}
	return dao.Announcements.Create(ctx, map[string]interface{}{
		"id":         idgen.New(),
		"title":      title,
		"content":    content,
		"type":       atype,
		"created_by": createdBy,
	})
}

func (s *sAnnouncement) ListAdmin(ctx context.Context, page, size int) ([]*v1.AnnouncementItem, int, error) {
	rows, total, err := dao.Announcements.List(ctx, page, size)
	if err != nil {
		return nil, 0, err
	}
	list := make([]*v1.AnnouncementItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, rowToItem(row, false))
	}
	return list, total, nil
}

func (s *sAnnouncement) Update(ctx context.Context, id int64, title, content, atype string) error {
	data := g.Map{"title": title, "content": content}
	if atype != "" {
		data["type"] = atype
	}
	return dao.Announcements.Update(ctx, id, data)
}

func (s *sAnnouncement) Delete(ctx context.Context, id int64) error {
	return dao.Announcements.SoftDelete(ctx, id)
}

func (s *sAnnouncement) ListForUser(ctx context.Context, userId int64, page, size int) ([]*v1.AnnouncementItem, int, int, error) {
	cursor, err := dao.Announcements.GetCursor(ctx, userId)
	if err != nil {
		return nil, 0, 0, err
	}

	rows, total, err := dao.Announcements.List(ctx, page, size)
	if err != nil {
		return nil, 0, 0, err
	}

	list := make([]*v1.AnnouncementItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, rowToItem(row, isUnreadForCursor(row.CreatedAt, cursor)))
	}

	totalUnread, err := dao.Announcements.UnreadCount(ctx, userId)
	if err != nil {
		return nil, 0, 0, err
	}

	return list, total, totalUnread, nil
}

func (s *sAnnouncement) MarkRead(ctx context.Context, userId int64) error {
	return dao.Announcements.UpsertCursor(ctx, userId, time.Now())
}

func (s *sAnnouncement) UnreadCount(ctx context.Context, userId int64) (int, error) {
	return dao.Announcements.UnreadCount(ctx, userId)
}

// ── helpers ──────────────────────────────────────────────────────────────────

func rowToItem(row dao.AnnouncementRow, unread bool) *v1.AnnouncementItem {
	atype := row.Type
	if atype == "" {
		atype = "info"
	}
	return &v1.AnnouncementItem{
		Id:        row.Id,
		Title:     row.Title,
		Content:   row.Content,
		Type:      atype,
		CreatedAt: row.CreatedAt,
		Unread:    unread,
	}
}

func isUnreadForCursor(createdAtStr string, cursor *time.Time) bool {
	if cursor == nil {
		return true
	}
	if createdAtStr == "" {
		return false
	}
	t, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		t, err = time.Parse("2006-01-02 15:04:05.000", createdAtStr)
	}
	if err != nil {
		return false
	}
	return t.After(*cursor)
}
