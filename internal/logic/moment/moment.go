package moment

import (
	"context"
	"math"

	momentv1 "github.com/nuxtblog/nuxtblog/api/moment/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/model/entity"
	"github.com/nuxtblog/nuxtblog/internal/search"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/nuxtblog/nuxtblog/internal/util/idgen"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sMoment struct{}

func New() *sMoment { return &sMoment{} }

func init() {
	service.RegisterMoment(New())
}

func (s *sMoment) Create(ctx context.Context, req *momentv1.MomentCreateReq) (int64, error) {
	uid, _ := middleware.GetCurrentUserID(ctx)
	visibility := req.Visibility
	if visibility == 0 {
		visibility = 1
	}
	m := &entity.Moments{
		Id:         int(idgen.New()),
		AuthorId:   int(uid),
		Content:    req.Content,
		Visibility: visibility,
	}
	id, err := dao.Moments.Create(ctx, m)
	if err != nil {
		return 0, err
	}

	// Attach media if provided
	if len(req.MediaIds) > 0 {
		rows := make([]g.Map, 0, len(req.MediaIds))
		for i, mid := range req.MediaIds {
			rows = append(rows, g.Map{
				"moment_id":  id,
				"media_id":   mid,
				"sort_order": i,
			})
		}
		_, _ = dao.MomentMedia.Ctx(ctx).Data(rows).Insert()
	}

	// Insert initial moment_stats
	_, _ = dao.MomentStats.Ctx(ctx).Data(g.Map{
		"moment_id":     id,
		"view_count":    0,
		"like_count":    0,
		"comment_count": 0,
	}).Insert()

	return id, nil
}

func (s *sMoment) Update(ctx context.Context, req *momentv1.MomentUpdateReq) error {
	uid, _ := middleware.GetCurrentUserID(ctx)
	role := middleware.GetCurrentUserRole(ctx)

	// Check ownership
	type ownerRow struct {
		AuthorId int64 `orm:"author_id"`
	}
	var owner ownerRow
	if err := dao.Moments.Ctx(ctx).Where("id", req.Id).WhereNull("deleted_at").Scan(&owner); err != nil {
		return err
	}
	if owner.AuthorId == 0 {
		return gerror.NewCode(gcode.CodeNotFound, "moment not found")
	}
	if owner.AuthorId != uid && role < middleware.RoleAdmin {
		return gerror.NewCode(gcode.CodeNotAuthorized, "permission denied")
	}

	data := g.Map{}
	if req.Content != nil {
		data["content"] = *req.Content
	}
	if req.Visibility != nil {
		data["visibility"] = *req.Visibility
	}
	if len(data) > 0 {
		if _, err := dao.Moments.Ctx(ctx).Where("id", req.Id).WhereNull("deleted_at").Data(data).Update(); err != nil {
			return err
		}
	}

	// Replace media if provided
	if req.MediaIds != nil {
		if _, err := dao.MomentMedia.Ctx(ctx).Where("moment_id", req.Id).Delete(); err != nil {
			return err
		}
		if len(*req.MediaIds) > 0 {
			rows := make([]g.Map, 0, len(*req.MediaIds))
			for i, mid := range *req.MediaIds {
				rows = append(rows, g.Map{
					"moment_id":  req.Id,
					"media_id":   mid,
					"sort_order": i,
				})
			}
			if _, err := dao.MomentMedia.Ctx(ctx).Data(rows).Insert(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *sMoment) Delete(ctx context.Context, id int64) error {
	uid, _ := middleware.GetCurrentUserID(ctx)
	role := middleware.GetCurrentUserRole(ctx)

	type ownerRow struct {
		AuthorId int64 `orm:"author_id"`
	}
	var owner ownerRow
	if err := dao.Moments.Ctx(ctx).Where("id", id).WhereNull("deleted_at").Scan(&owner); err != nil {
		return err
	}
	if owner.AuthorId == 0 {
		return gerror.NewCode(gcode.CodeNotFound, "moment not found")
	}
	if owner.AuthorId != uid && role < middleware.RoleAdmin {
		return gerror.NewCode(gcode.CodeNotAuthorized, "permission denied")
	}

	_, err := dao.Moments.Ctx(ctx).
		Where("id", id).WhereNull("deleted_at").
		Data(g.Map{"deleted_at": gtime.Now()}).Update()
	return err
}

func (s *sMoment) GetById(ctx context.Context, id int64) (*momentv1.MomentItem, error) {
	m, err := dao.Moments.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}
	return s.buildMomentItem(ctx, m)
}

func (s *sMoment) buildMomentItem(ctx context.Context, m *entity.Moments) (*momentv1.MomentItem, error) {
	item := &momentv1.MomentItem{
		Id:         int64(m.Id),
		AuthorId:   int64(m.AuthorId),
		Content:    m.Content,
		Visibility: m.Visibility,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}

	// Load stats
	type statsRow struct {
		ViewCount    int64 `orm:"view_count"`
		LikeCount    int64 `orm:"like_count"`
		CommentCount int64 `orm:"comment_count"`
	}
	var stats statsRow
	if err := dao.MomentStats.Ctx(ctx).Where("moment_id", m.Id).Scan(&stats); err == nil {
		item.Stats = &momentv1.MomentStatsItem{
			ViewCount:    stats.ViewCount,
			LikeCount:    stats.LikeCount,
			CommentCount: stats.CommentCount,
		}
	}

	// Load media
	type mediaRow struct {
		MediaId   int64  `orm:"media_id"`
		SortOrder int    `orm:"sort_order"`
	}
	var mediaRows []mediaRow
	if err := dao.MomentMedia.Ctx(ctx).Where("moment_id", m.Id).OrderAsc("sort_order").Scan(&mediaRows); err == nil && len(mediaRows) > 0 {
		mediaIds := make([]int64, 0, len(mediaRows))
		for _, mr := range mediaRows {
			mediaIds = append(mediaIds, mr.MediaId)
		}
		type mediaDet struct {
			Id       int64  `orm:"id"`
			CdnUrl   string `orm:"cdn_url"`
			MimeType string `orm:"mime_type"`
			Width    int    `orm:"width"`
			Height   int    `orm:"height"`
		}
		var medias []mediaDet
		if err2 := dao.Medias.Ctx(ctx).WhereIn("id", mediaIds).WhereNull("deleted_at").Scan(&medias); err2 == nil {
			mediaMap := make(map[int64]*mediaDet)
			for i := range medias {
				mediaMap[medias[i].Id] = &medias[i]
			}
			for _, mr := range mediaRows {
				if md, ok := mediaMap[mr.MediaId]; ok {
					item.Media = append(item.Media, &momentv1.MomentMediaItem{
						Id:       md.Id,
						Url:      md.CdnUrl,
						MimeType: md.MimeType,
						Width:    md.Width,
						Height:   md.Height,
					})
				}
			}
		}
	}

	// Load author
	type authorRow struct {
		Id          int64  `orm:"id"`
		Username    string `orm:"username"`
		DisplayName string `orm:"display_name"`
	}
	var author authorRow
	if err := dao.Users.Ctx(ctx).Where("id", m.AuthorId).Scan(&author); err == nil && author.Id > 0 {
		item.Author = &momentv1.MomentAuthorItem{
			Id:       author.Id,
			Username: author.Username,
			Nickname: author.DisplayName,
		}
	}

	return item, nil
}

func (s *sMoment) GetList(ctx context.Context, req *momentv1.MomentGetListReq) (*momentv1.MomentGetListRes, error) {
	uid, _ := middleware.GetCurrentUserID(ctx)
	role := middleware.GetCurrentUserRole(ctx)

	m := dao.Moments.Ctx(ctx).WhereNull("deleted_at")
	if req.AuthorId != nil {
		m = m.Where("author_id", *req.AuthorId)
	}
	if req.Visibility != nil {
		m = m.Where("visibility", *req.Visibility)
	} else {
		// Non-admins only see public moments unless they are the author
		if role < middleware.RoleAdmin {
			if uid > 0 {
				m = m.Where(m.Builder().Where("visibility", 1).WhereOr("author_id", uid))
			} else {
				m = m.Where("visibility", 1)
			}
		}
	}
	if req.Keyword != nil && *req.Keyword != "" {
		sr, err := search.Default(ctx).Search(ctx, search.ContentMoment, *req.Keyword)
		if err != nil {
			return nil, err
		}
		if len(sr.IDs) == 0 {
			return &momentv1.MomentGetListRes{Data: []*momentv1.MomentItem{}, Page: req.Page, PageSize: req.PageSize}, nil
		}
		m = m.WhereIn("id", sr.IDs)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	type momentRow struct {
		Id         int64       `orm:"id"`
		AuthorId   int64       `orm:"author_id"`
		Content    string      `orm:"content"`
		Visibility int         `orm:"visibility"`
		CreatedAt  *gtime.Time `orm:"created_at"`
		UpdatedAt  *gtime.Time `orm:"updated_at"`
	}
	var rows []momentRow
	if total > 0 {
		if err = m.Page(req.Page, req.PageSize).OrderDesc("created_at").Scan(&rows); err != nil {
			return nil, err
		}
	}

	data := make([]*momentv1.MomentItem, 0, len(rows))
	for _, r := range rows {
		ent := &entity.Moments{
			Id:         int(r.Id),
			AuthorId:   int(r.AuthorId),
			Content:    r.Content,
			Visibility: r.Visibility,
			CreatedAt:  r.CreatedAt,
			UpdatedAt:  r.UpdatedAt,
		}
		item, err := s.buildMomentItem(ctx, ent)
		if err != nil {
			continue
		}
		data = append(data, item)
	}

	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	return &momentv1.MomentGetListRes{
		Data:       data,
		Total:      total,
		Page:       req.Page,
		PageSize:   pageSize,
		TotalPages: int(math.Ceil(float64(total) / float64(pageSize))),
	}, nil
}

func (s *sMoment) IncrementView(ctx context.Context, id int64) error {
	_, err := dao.MomentStats.DB().Exec(
		ctx,
		`INSERT INTO moment_stats (moment_id, view_count, like_count, comment_count)
		 VALUES (?, 1, 0, 0)
		 ON CONFLICT(moment_id) DO UPDATE SET view_count = view_count + 1`,
		id,
	)
	if err != nil {
		g.Log().Warningf(ctx, "[moment] increment view error: %v", err)
	}
	return nil
}
