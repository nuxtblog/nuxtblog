package message

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/message/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/nuxtblog/nuxtblog/internal/util/idgen"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sMessage struct{}

func New() service.IMessage { return &sMessage{} }
func init()                 { service.RegisterMessage(New()) }

func requireLogin(ctx context.Context) (int64, error) {
	id, _ := middleware.GetCurrentUserID(ctx)
	if id == 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, g.I18n().T(ctx, "error.unauthorized"))
	}
	return id, nil
}

// getOrCreateConversation returns conversation id, ensuring user_a < user_b for uniqueness.
func getOrCreateConversation(ctx context.Context, me, other int64) (int64, error) {
	a, b := me, other
	if a > b {
		a, b = b, a
	}
	type convRow struct {
		Id int64 `orm:"id"`
	}
	var row convRow
	_ = dao.Conversations.Ctx(ctx).
		Where("user_a", a).Where("user_b", b).Fields("id").Scan(&row)
	if row.Id > 0 {
		return row.Id, nil
	}
	// Use a random ID to avoid sequential enumeration of conversation counts.
	newId := idgen.New()
	_, err := dao.Conversations.Ctx(ctx).Data(g.Map{
		"id": newId, "user_a": a, "user_b": b,
	}).Insert()
	return newId, err
}

func (s *sMessage) Send(ctx context.Context, req *v1.MessageSendReq) (*v1.MessageSendRes, error) {
	me, err := requireLogin(ctx)
	if err != nil {
		return nil, err
	}
	if me == req.ToUserId {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, g.I18n().T(ctx, "message.cannot_send_to_self"))
	}

	convId, err := getOrCreateConversation(ctx, me, req.ToUserId)
	if err != nil {
		return nil, err
	}

	// Use a random ID to avoid sequential enumeration of message counts.
	msgId := idgen.New()
	_, err = dao.Messages.Ctx(ctx).Data(g.Map{
		"id":              msgId,
		"conversation_id": convId,
		"sender_id":       me,
		"content":         req.Content,
	}).Insert()
	if err != nil {
		return nil, err
	}

	// Update conversation last message + unread count for the other user
	a, b := me, req.ToUserId
	if a > b {
		a, b = b, a
	}
	isRecipientA := req.ToUserId == a
	convM := dao.Conversations.Ctx(ctx).Where("id", convId).Data(g.Map{
		"last_msg":    req.Content,
		"last_msg_at": gtime.Now(),
	})
	if _, err2 := convM.Update(); err2 == nil {
		if isRecipientA {
			_, _ = dao.Conversations.Ctx(ctx).Where("id", convId).Increment("unread_a", 1)
		} else {
			_, _ = dao.Conversations.Ctx(ctx).Where("id", convId).Increment("unread_b", 1)
		}
	}

	return &v1.MessageSendRes{Id: msgId, ConversationId: convId}, nil
}

func (s *sMessage) ConversationList(ctx context.Context, req *v1.ConversationListReq) (*v1.ConversationListRes, error) {
	me, err := requireLogin(ctx)
	if err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.Size
	type row struct {
		Id          int64  `orm:"id"`
		OtherUserId int64  `orm:"other_user_id"`
		OtherName   string `orm:"other_name"`
		OtherAvatar string `orm:"other_avatar"`
		LastMsg     string `orm:"last_msg"`
		LastMsgAt   string `orm:"last_msg_at"`
		UnreadCount int    `orm:"unread_count"`
	}

	var rows []row
	_ = dao.Conversations.DB().Ctx(ctx).Raw(`
		SELECT c.id,
		  CASE WHEN c.user_a = ? THEN c.user_b ELSE c.user_a END AS other_user_id,
		  COALESCE(u.display_name, u.username, '') AS other_name,
		  COALESCE(m.cdn_url, '') AS other_avatar,
		  c.last_msg,
		  COALESCE(c.last_msg_at, c.created_at) AS last_msg_at,
		  CASE WHEN c.user_a = ? THEN c.unread_a ELSE c.unread_b END AS unread_count
		FROM conversations c
		LEFT JOIN users u ON u.id = CASE WHEN c.user_a = ? THEN c.user_b ELSE c.user_a END
		LEFT JOIN medias m ON m.id = u.avatar_id
		WHERE c.user_a = ? OR c.user_b = ?
		ORDER BY last_msg_at DESC NULLS LAST
		LIMIT ? OFFSET ?
	`, me, me, me, me, me, req.Size, offset).Scan(&rows)

	var total int
	_ = dao.Conversations.DB().Ctx(ctx).Raw(
		`SELECT COUNT(*) FROM conversations WHERE user_a=? OR user_b=?`, me, me,
	).Scan(&total)

	var totalUnread int
	_ = dao.Conversations.DB().Ctx(ctx).Raw(`
		SELECT COALESCE(SUM(CASE WHEN user_a=? THEN unread_a ELSE unread_b END), 0)
		FROM conversations WHERE user_a=? OR user_b=?
	`, me, me, me).Scan(&totalUnread)

	items := make([]v1.ConversationItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, v1.ConversationItem{
			Id: r.Id, OtherUserId: r.OtherUserId,
			OtherName: r.OtherName, OtherAvatar: r.OtherAvatar,
			LastMsg: r.LastMsg, LastMsgAt: r.LastMsgAt, UnreadCount: r.UnreadCount,
		})
	}
	return &v1.ConversationListRes{Items: items, Total: total, TotalUnread: totalUnread}, nil
}

func (s *sMessage) MessageList(ctx context.Context, req *v1.MessageListReq) (*v1.MessageListRes, error) {
	me, err := requireLogin(ctx)
	if err != nil {
		return nil, err
	}

	a, b := me, req.ToUserId
	if a > b {
		a, b = b, a
	}
	type convRow struct {
		Id int64 `orm:"id"`
	}
	var conv convRow
	_ = dao.Conversations.Ctx(ctx).
		Where("user_a", a).Where("user_b", b).Fields("id").Scan(&conv)

	if conv.Id == 0 {
		return &v1.MessageListRes{Items: []v1.MessageItem{}, ConversationId: 0}, nil
	}

	type row struct {
		Id        int64  `orm:"id"`
		SenderId  int64  `orm:"sender_id"`
		Content   string `orm:"content"`
		IsRead    bool   `orm:"is_read"`
		CreatedAt string `orm:"created_at"`
	}

	q := dao.Messages.Ctx(ctx).
		Where("conversation_id", conv.Id).
		OrderDesc("id").
		Limit(req.Size + 1)
	if req.BeforeId > 0 {
		q = q.WhereLT("id", req.BeforeId)
	}
	var rows []row
	_ = q.Scan(&rows)

	hasMore := len(rows) > req.Size
	if hasMore {
		rows = rows[:req.Size]
	}

	// Mark as read
	_, _ = dao.Messages.Ctx(ctx).
		Where("conversation_id", conv.Id).
		Where("sender_id", req.ToUserId).
		Where("is_read", 0).
		Data(g.Map{"is_read": 1}).Update()

	// Reset my unread counter (use two separate Update calls to avoid field name injection)
	if me == b {
		_, _ = dao.Conversations.Ctx(ctx).Where("id", conv.Id).Data(g.Map{"unread_b": 0}).Update()
	} else {
		_, _ = dao.Conversations.Ctx(ctx).Where("id", conv.Id).Data(g.Map{"unread_a": 0}).Update()
	}

	items := make([]v1.MessageItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, v1.MessageItem{
			Id: r.Id, SenderId: r.SenderId, Content: r.Content,
			IsRead: r.IsRead, CreatedAt: r.CreatedAt,
		})
	}
	// Reverse to chronological order
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}
	return &v1.MessageListRes{Items: items, ConversationId: conv.Id, HasMore: hasMore}, nil
}

func (s *sMessage) UnreadCount(ctx context.Context) (int, error) {
	me, err := requireLogin(ctx)
	if err != nil {
		return 0, nil // silently return 0 for unauthenticated
	}
	var count int
	_ = dao.Conversations.DB().Ctx(ctx).Raw(`
		SELECT COALESCE(SUM(CASE WHEN user_a=? THEN unread_a ELSE unread_b END), 0)
		FROM conversations WHERE user_a=? OR user_b=?
	`, me, me, me).Scan(&count)
	return count, nil
}
