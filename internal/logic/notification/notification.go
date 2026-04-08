package notification

import (
	"context"
	"math"
	"time"

	v1 "github.com/nuxtblog/nuxtblog/api/notification/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/notify"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/nuxtblog/nuxtblog/internal/util/idgen"

	"github.com/gogf/gf/v2/frame/g"
)

type sNotification struct{}

func New() service.INotification { return &sNotification{} }

func init() {
	service.RegisterNotification(New())
}


func (s *sNotification) Create(ctx context.Context, notifType, subType string, actorId *int64, actorName, actorAvatar string, userId int64, objectType string, objectId *int64, objectTitle, objectLink, content string) error {
	if len(content) > 200 {
		content = content[:200]
	}
	data := g.Map{
		"id":           idgen.New(),
		"user_id":      userId,
		"type":         notifType,
		"sub_type":     subType,
		"actor_id":     actorId,
		"actor_name":   actorName,
		"actor_avatar": actorAvatar,
		"object_type":  objectType,
		"object_id":    objectId,
		"object_title": objectTitle,
		"object_link":  objectLink,
		"content":      content,
		"is_read":      0,
	}
	_, err := dao.Notifications.Create(ctx, data)
	if err != nil {
		return err
	}
	// Fan-out to extra channels (email, SMS, …) asynchronously.
	go notify.Dispatch(context.Background(), notify.Message{
		UserID:      userId,
		Type:        notifType,
		SubType:     subType,
		ActorID:     actorId,
		ActorName:   actorName,
		ActorAvatar: actorAvatar,
		ObjectType:  objectType,
		ObjectID:    objectId,
		ObjectTitle: objectTitle,
		ObjectLink:  objectLink,
		Content:     content,
	})
	// Async: prune read notifications older than 90 days to keep the table lean
	go func() {
		tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cutoff := time.Now().AddDate(0, 0, -90).Format("2006-01-02 15:04:05")
		_, _ = dao.Notifications.Ctx(tctx).
			Where("user_id", userId).
			Where("is_read", 1).
			WhereLT("created_at", cutoff).
			Delete()
	}()
	return nil
}

func (s *sNotification) List(ctx context.Context, req *v1.NotificationListReq) (*v1.NotificationListRes, error) {
	m := dao.Notifications.Ctx(ctx).
		Where("user_id", req.UserId).
		WhereNull("deleted_at")

	switch req.Filter {
	case "unread":
		m = m.Where("is_read", 0)
	case "interaction":
		m = m.WhereNot("type", "system")
	case "system":
		m = m.Where("type", "system")
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	notifUnread, err := dao.Notifications.UnreadCount(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	announcUnread, _ := dao.Announcements.UnreadCount(ctx, req.UserId)
	unread := notifUnread + announcUnread

	type Row struct {
		Id          int64  `orm:"id"`
		Type        string `orm:"type"`
		SubType     string `orm:"sub_type"`
		ActorId     *int64 `orm:"actor_id"`
		ActorName   string `orm:"actor_name"`
		ActorAvatar string `orm:"actor_avatar"`
		ObjectType  string `orm:"object_type"`
		ObjectId    *int64 `orm:"object_id"`
		ObjectTitle string `orm:"object_title"`
		ObjectLink  string `orm:"object_link"`
		Content     string `orm:"content"`
		IsRead      int    `orm:"is_read"`
		CreatedAt   string `orm:"created_at"`
	}

	var rows []Row
	if total > 0 {
		err = m.Page(req.Page, req.Size).OrderDesc("created_at").Scan(&rows)
		if err != nil {
			return nil, err
		}
	}

	list := make([]*v1.NotificationItem, len(rows))
	for i, row := range rows {
		item := &v1.NotificationItem{
			Id:           row.Id,
			Type:         row.Type,
			SubType:      row.SubType,
			UserName:     row.ActorName,
			Avatar:       row.ActorAvatar,
			Content:      row.Content,
			RelatedTitle: row.ObjectTitle,
			RelatedLink:  row.ObjectLink,
			Read:         row.IsRead == 1,
			CreatedAt:    row.CreatedAt,
		}
		if row.Type == "system" {
			item.Title = g.I18n().T(ctx, "notification.system_"+row.SubType)
		} else {
			item.Action = g.I18n().T(ctx, "notification.action_"+row.Type)
		}
		list[i] = item
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Size)))

	return &v1.NotificationListRes{
		List:       list,
		Total:      total,
		Page:       req.Page,
		Size:       req.Size,
		TotalPages: totalPages,
		Unread:     unread,
	}, nil
}

func (s *sNotification) UnreadCount(ctx context.Context, userId int64) (int, error) {
	notifCount, err := dao.Notifications.UnreadCount(ctx, userId)
	if err != nil {
		return 0, err
	}
	announcCount, err := dao.Announcements.UnreadCount(ctx, userId)
	if err != nil {
		return 0, err
	}
	return notifCount + announcCount, nil
}

func (s *sNotification) MarkRead(ctx context.Context, id int64) error {
	return dao.Notifications.MarkRead(ctx, id)
}

func (s *sNotification) MarkAllRead(ctx context.Context, userId int64) error {
	return dao.Notifications.MarkAllRead(ctx, userId)
}

func (s *sNotification) Delete(ctx context.Context, id int64) error {
	return dao.Notifications.SoftDelete(ctx, id)
}

func (s *sNotification) Clear(ctx context.Context, userId int64, filter string) error {
	return dao.Notifications.ClearByFilter(ctx, userId, filter)
}
