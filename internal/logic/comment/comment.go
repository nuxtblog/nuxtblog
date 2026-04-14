package comment

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/nuxtblog/nuxtblog/api/comment/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/event"
	"github.com/nuxtblog/nuxtblog/internal/event/payload"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	eng "github.com/nuxtblog/nuxtblog/internal/pluginsys"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/nuxtblog/nuxtblog/internal/util/idgen"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type sComment struct{}

func New() service.IComment { return &sComment{} }

func init() {
	service.RegisterComment(New())
}

var commentStatusIntMap = map[string]int{
	"1": 1, "pending": 1,
	"2": 2, "approved": 2,
	"3": 3, "spam": 3,
	"4": 4, "trash": 4,
}

var commentStatusStrMap = map[int]string{
	1: "pending",
	2: "approved",
	3: "spam",
	4: "trash",
}

func (s *sComment) AdminGetList(ctx context.Context, req *v1.CommentAdminGetListReq) (*v1.CommentAdminGetListRes, error) {
	m := dao.Comments.Ctx(ctx).WhereNull("deleted_at")
	if req.Status != nil && *req.Status != "" && *req.Status != "all" {
		if si, ok := commentStatusIntMap[*req.Status]; ok {
			m = m.Where("status", si)
		}
	}
	if req.ObjectType != nil && *req.ObjectType != "" {
		m = m.Where("object_type", *req.ObjectType)
	}
	if req.Keyword != nil && *req.Keyword != "" {
		kw := fmt.Sprintf("%%%s%%", *req.Keyword)
		m = m.WhereOrLike("content", kw).WhereOrLike("author_name", kw)
	}
	m = m.WhereNull("parent_id")

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	type CommentRow struct {
		Id          int64  `orm:"id"`
		ObjectId    int64  `orm:"object_id"`
		ObjectType  string `orm:"object_type"`
		ParentId    *int64 `orm:"parent_id"`
		UserId      *int64 `orm:"user_id"`
		AuthorName  string `orm:"author_name"`
		AuthorEmail string `orm:"author_email"`
		Content     string `orm:"content"`
		Status      int    `orm:"status"`
		Ip          string `orm:"ip"`
		UserAgent   string `orm:"user_agent"`
		CreatedAt   string `orm:"created_at"`
	}

	var rows []CommentRow
	if total > 0 {
		err = m.Page(req.Page, req.Size).OrderDesc("created_at").Scan(&rows)
		if err != nil {
			return nil, err
		}
	}

	parentIds := make([]int64, 0, len(rows))
	for _, row := range rows {
		parentIds = append(parentIds, row.Id)
	}

	userIdSet := map[int64]bool{}
	for _, row := range rows {
		if row.UserId != nil {
			userIdSet[*row.UserId] = true
		}
	}

	childrenMap := map[int64][]*CommentRow{}
	if len(parentIds) > 0 {
		var childRows []CommentRow
		_ = dao.Comments.Ctx(ctx).
			WhereIn("parent_id", parentIds).
			WhereNull("deleted_at").
			OrderAsc("created_at").
			Scan(&childRows)
		for i := range childRows {
			pid := *childRows[i].ParentId
			childrenMap[pid] = append(childrenMap[pid], &childRows[i])
			if childRows[i].UserId != nil {
				userIdSet[*childRows[i].UserId] = true
			}
		}
	}

	avatarMap := map[int64]string{}
	if len(userIdSet) > 0 {
		userIds := make([]int64, 0, len(userIdSet))
		for id := range userIdSet {
			userIds = append(userIds, id)
		}
		type AvatarRow struct {
			Id       int64 `orm:"id"`
			AvatarId int64 `orm:"avatar_id"`
		}
		var avatarRows []AvatarRow
		_ = dao.Users.Ctx(ctx).Fields("id, avatar_id").WhereIn("id", userIds).Scan(&avatarRows)
		for _, ar := range avatarRows {
			if ar.AvatarId > 0 {
				val, err := dao.Medias.Ctx(ctx).Fields("cdn_url").Where("id", ar.AvatarId).Value()
				if err == nil && !val.IsEmpty() {
					avatarMap[ar.Id] = val.String()
				}
			}
		}
	}

	buildItem := func(row *CommentRow) *v1.CommentAdminItem {
		item := &v1.CommentAdminItem{
			CommentId: row.Id,
			PostId:    row.ObjectId,
			Content:   row.Content,
			ParentId:  row.ParentId,
			Status:    commentStatusStrMap[row.Status],
			IpAddress: row.Ip,
			UserAgent: row.UserAgent,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.CreatedAt,
			Author: &v1.CommentAdminAuthor{
				Name:  row.AuthorName,
				Email: row.AuthorEmail,
			},
		}
		if row.UserId != nil {
			item.Author.Id = row.UserId
			item.Author.Avatar = avatarMap[*row.UserId]
		}
		return item
	}

	list := make([]*v1.CommentAdminItem, len(rows))
	for i, row := range rows {
		item := buildItem(&row)
		if children, ok := childrenMap[row.Id]; ok {
			for _, c := range children {
				item.Children = append(item.Children, buildItem(c))
			}
		}
		list[i] = item
	}

	return &v1.CommentAdminGetListRes{
		List:  list,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}

func (s *sComment) Create(ctx context.Context, req *v1.CommentCreateReq) (*v1.CommentCreateRes, error) {
	// Determine initial status based on comment_moderation setting
	initialStatus := int(v1.CommentReviewApproved)
	if val, err2 := dao.Options.Ctx(ctx).Fields("value").Where("key", "comment_moderation").Value(); err2 == nil {
		v := val.String()
		if v == "true" || v == "1" {
			initialStatus = int(v1.CommentReviewPending)
		}
	}

	// Run plugin filter:comment.create — allows plugins to modify or reject comments
	authorName := req.AuthorName
	authorEmail := req.AuthorEmail
	content := req.Content
	filtered, filterErr := eng.Filter(ctx, eng.FilterCommentCreate, map[string]any{
		"content":      content,
		"author_name":  authorName,
		"author_email": authorEmail,
	})
	if filterErr != nil {
		return nil, filterErr
	}
	if v, ok := filtered["content"].(string); ok {
		content = v
	}
	if v, ok := filtered["author_name"].(string); ok {
		authorName = v
	}
	if v, ok := filtered["author_email"].(string); ok {
		authorEmail = v
	}

	row := g.Map{
		"id":           idgen.New(),
		"object_id":    req.ObjectId,
		"object_type":  req.ObjectType,
		"parent_id":    req.ParentId,
		"author_name":  authorName,
		"author_email": authorEmail,
		"content":      content,
		"status":       initialStatus,
	}
	// Associate comment with the logged-in user if authenticated
	if r := ghttp.RequestFromCtx(ctx); r != nil {
		if uid := r.GetCtxVar("user_id").Int64(); uid > 0 {
			row["user_id"] = uid
		}
	}
	result, err := dao.Comments.Ctx(ctx).Insert(row)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()

	// Emit comment.created — notification delivery is handled by the listener.
	if req.ObjectType == "post" {
		type PostRow struct {
			AuthorId int64  `orm:"author_id"`
			Title    string `orm:"title"`
			Slug     string `orm:"slug"`
		}
		var post PostRow
		_ = dao.Posts.Ctx(ctx).Fields("author_id, title, slug").Where("id", req.ObjectId).Scan(&post)

		var authorID int64
		if r := ghttp.RequestFromCtx(ctx); r != nil {
			authorID = r.GetCtxVar("user_id").Int64()
		}

		p := payload.CommentCreated{
			CommentID:    id,
			Status:       initialStatus,
			ObjectType:   req.ObjectType,
			ObjectID:     req.ObjectId,
			ObjectTitle:  post.Title,
			ObjectSlug:   post.Slug,
			PostAuthorID: post.AuthorId,
			ParentID:     req.ParentId,
			AuthorID:     authorID,
			AuthorName:   authorName,
			AuthorEmail:  req.AuthorEmail,
			Content:      content,
		}
		if req.ParentId != nil {
			type ParentRow struct {
				UserId *int64 `orm:"user_id"`
			}
			var parent ParentRow
			_ = dao.Comments.Ctx(ctx).Fields("user_id").Where("id", *req.ParentId).Scan(&parent)
			if parent.UserId != nil && *parent.UserId != post.AuthorId {
				p.ParentAuthorID = *parent.UserId
			}
		}
		_ = event.Emit(ctx, event.CommentCreated, p)
	}

	return &v1.CommentCreateRes{Id: id, Status: v1.CommentReviewStatus(initialStatus)}, nil
}

func (s *sComment) Delete(ctx context.Context, id int64) error {
	role := middleware.GetCurrentUserRole(ctx)
	if !service.Permission().Can(ctx, role, "manage_comments") {
		return gerror.NewCode(gcode.CodeNotAuthorized, "permission denied: manage_comments")
	}

	type commentRow struct {
		ObjectType string `orm:"object_type"`
		ObjectID   int64  `orm:"object_id"`
	}
	var crow commentRow
	_ = dao.Comments.Ctx(ctx).Fields("object_type, object_id").Where("id", id).Scan(&crow)

	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := dao.Comments.Ctx(ctx).
		Where("id", id).
		WhereNull("deleted_at").
		Update(g.Map{"deleted_at": now})
	if err != nil {
		return err
	}
	_, _ = dao.Comments.Ctx(ctx).
		Where("parent_id", id).
		WhereNull("deleted_at").
		Update(g.Map{"deleted_at": now})
	_ = event.Emit(ctx, event.CommentDeleted, payload.CommentDeleted{
		CommentID: id, ObjectType: crow.ObjectType, ObjectID: crow.ObjectID,
	})
	return nil
}

func (s *sComment) GetList(ctx context.Context, req *v1.CommentGetListReq) (*v1.CommentGetListRes, error) {
	m := dao.Comments.Ctx(ctx).
		WhereNull("deleted_at").
		Where("object_id", req.ObjectId).
		Where("object_type", req.ObjectType).
		Where("status", int(v1.CommentReviewApproved)).
		WhereNull("parent_id")

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	res := &v1.CommentGetListRes{
		List:  []*v1.CommentItem{},
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}
	if total == 0 {
		return res, nil
	}

	type CommentRow struct {
		Id          int64  `orm:"id"`
		ObjectId    int64  `orm:"object_id"`
		ObjectType  string `orm:"object_type"`
		ParentId    *int64 `orm:"parent_id"`
		UserId      *int64 `orm:"user_id"`
		AuthorName  string `orm:"author_name"`
		AuthorEmail string `orm:"author_email"`
		Content     string `orm:"content"`
		Status      int    `orm:"status"`
	}

	var rows []CommentRow
	err = m.Page(req.Page, req.Size).OrderDesc("created_at").Scan(&rows)
	if err != nil {
		return nil, err
	}

	// Load replies
	parentIds := make([]int64, 0, len(rows))
	for _, r := range rows {
		parentIds = append(parentIds, r.Id)
	}
	repliesMap := map[int64][]*v1.CommentItem{}
	if len(parentIds) > 0 {
		var replyRows []CommentRow
		_ = dao.Comments.Ctx(ctx).
			WhereIn("parent_id", parentIds).
			WhereNull("deleted_at").
			Where("status", int(v1.CommentReviewApproved)).
			OrderAsc("created_at").
			Scan(&replyRows)
		for _, r := range replyRows {
			repliesMap[*r.ParentId] = append(repliesMap[*r.ParentId], &v1.CommentItem{
				Id:         r.Id,
				ObjectId:   r.ObjectId,
				ObjectType: r.ObjectType,
				ParentId:   r.ParentId,
				UserId:     r.UserId,
				AuthorName: r.AuthorName,
				Content:    r.Content,
				Status:     v1.CommentReviewStatus(r.Status),
			})
		}
	}

	for _, row := range rows {
		item := &v1.CommentItem{
			Id:         row.Id,
			ObjectId:   row.ObjectId,
			ObjectType: row.ObjectType,
			ParentId:   row.ParentId,
			UserId:     row.UserId,
			AuthorName: row.AuthorName,
			Content:    row.Content,
			Status:     v1.CommentReviewStatus(row.Status),
			Replies:    repliesMap[row.Id],
		}
		res.List = append(res.List, item)
	}
	return res, nil
}

func (s *sComment) UpdateStatus(ctx context.Context, id int64, status v1.CommentReviewStatus) error {
	role := middleware.GetCurrentUserRole(ctx)
	if !service.Permission().Can(ctx, role, "moderate_comments") {
		return gerror.NewCode(gcode.CodeNotAuthorized, "permission denied: moderate_comments")
	}

	type snapRow struct {
		ObjectType string `orm:"object_type"`
		ObjectID   int64  `orm:"object_id"`
		Status     int    `orm:"status"`
	}
	var snap snapRow
	_ = dao.Comments.Ctx(ctx).Where("id", id).WhereNull("deleted_at").Scan(&snap)

	_, err := dao.Comments.Ctx(ctx).
		Where("id", id).
		WhereNull("deleted_at").
		Update(g.Map{"status": int(status)})
	if err != nil {
		return err
	}

	moderatorID, _ := middleware.GetCurrentUserID(ctx)
	_ = event.Emit(ctx, event.CommentStatusChanged, payload.CommentStatusChanged{
		CommentID:   id,
		ObjectType:  snap.ObjectType,
		ObjectID:    snap.ObjectID,
		OldStatus:   snap.Status,
		NewStatus:   int(status),
		ModeratorID: moderatorID,
	})
	// Convenience event: comment.approved
	if int(status) == 2 {
		_ = event.Emit(ctx, event.CommentApproved, payload.CommentApproved{
			CommentID:   id,
			ObjectType:  snap.ObjectType,
			ObjectID:    snap.ObjectID,
			ModeratorID: moderatorID,
		})
	}
	return nil
}

func (s *sComment) UpdateContent(ctx context.Context, id int64, content string) error {
	role := middleware.GetCurrentUserRole(ctx)
	if !service.Permission().Can(ctx, role, "manage_comments") {
		return gerror.NewCode(gcode.CodeNotAuthorized, "permission denied: manage_comments")
	}
	// Plugin filter — allows plugins to modify or reject comment edits
	filtered, ferr := eng.Filter(ctx, eng.FilterCommentUpdate, map[string]any{
		"id":      id,
		"content": content,
	})
	if ferr != nil {
		return ferr
	}
	if v, ok := filtered["content"].(string); ok {
		content = v
	}
	_, err := dao.Comments.Ctx(ctx).
		Where("id", id).
		WhereNull("deleted_at").
		Update(g.Map{"content": content})
	return err
}
