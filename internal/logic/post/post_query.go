package post

import (
	"context"
	"strconv"
	"strings"
	"time"

	v1 "github.com/nuxtblog/nuxtblog/api/post/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	plugin "github.com/nuxtblog/nuxtblog/internal/pluginsys"
	"github.com/nuxtblog/nuxtblog/internal/search"
	"github.com/nuxtblog/nuxtblog/internal/util/idgen"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// gtimeToRFC3339 formats a *gtime.Time as UTC RFC3339 ("2006-01-02T15:04:05Z").
// Returns "" for nil. Ensures frontend always receives timezone-aware ISO 8601 strings.
func gtimeToRFC3339(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.Time.UTC().Format(time.RFC3339)
}

// dbStrToRFC3339 parses a SQLite "YYYY-MM-DD HH:MM:SS" string (stored as server local time)
// and returns a UTC RFC3339 string. Falls back to the raw string on parse error.
func dbStrToRFC3339(s string) string {
	if s == "" {
		return ""
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	if err != nil {
		return s
	}
	return t.UTC().Format(time.RFC3339)
}

func parseIdList(s string) []int64 {
	var ids []int64
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if id, err := strconv.ParseInt(part, 10, 64); err == nil && id > 0 {
			ids = append(ids, id)
		}
	}
	return ids
}

// 枚举映射（供 logic 内部使用）
var postTypeStr = map[int]string{1: "post", 2: "page", 3: "custom"}
var postStatusStr = map[int]string{1: "draft", 2: "published", 3: "private", 4: "archived"}
var commentStatusStr = map[int]string{0: "closed", 1: "open"}

var statusIntMap = map[string]int{
	"1": 1, "draft": 1,
	"2": 2, "published": 2,
	"3": 3, "private": 3,
	"4": 4, "archived": 4,
	"5": 5, "trash": 5,
}

var typeIntMap = map[string]int{
	"1": 1, "post": 1,
	"2": 2, "page": 2,
	"3": 3, "custom": 3,
}

// ----------------------------------------------------------------
//  GetById  —  管理端用，返回原始字段 + Content
// ----------------------------------------------------------------

func (s *sPost) GetById(ctx context.Context, id int64) (item *v1.PostDetailItem, err error) {
	record, err := dao.Posts.GetById(ctx, id)
	if err != nil || record == nil {
		return nil, err
	}

	pi := &v1.PostItem{
		Id:            int64(record.Id),
		PostType:      v1.PostType(record.PostType),
		Status:        v1.PostStatus(record.Status),
		Title:         record.Title,
		Slug:          record.Slug,
		Excerpt:       record.Excerpt,
		AuthorId:      int64(record.AuthorId),
		CommentStatus: v1.CommentStatus(record.CommentStatus),
		Locale:        record.Locale,
		PublishedAt:   record.PublishedAt,
		CreatedAt:     record.CreatedAt,
		UpdatedAt:     record.UpdatedAt,
	}
	if record.FeaturedImgId != 0 {
		id64 := int64(record.FeaturedImgId)
		pi.FeaturedImgId = &id64
	}

	// Fetch stats for like_count
	type StatsRow struct {
		LikeCount int64 `orm:"like_count"`
	}
	var stats StatsRow
	if err2 := dao.PostStats.Ctx(ctx).Where("post_id", record.Id).Scan(&stats); err2 == nil && stats.LikeCount >= 0 {
		pi.Stats = &v1.PostStatsItem{LikeCount: stats.LikeCount}
	}

	return &v1.PostDetailItem{PostItem: *pi, Content: record.Content}, nil
}

// ----------------------------------------------------------------
//  GetBySlug  —  前台用，富化返回（作者、封面、统计、metas）
// ----------------------------------------------------------------

func (s *sPost) GetBySlug(ctx context.Context, slug string) (*v1.PostDetailEnrichedItem, error) {
	post, err := dao.Posts.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, nil
	}

	item := &v1.PostDetailEnrichedItem{
		Id:            int64(post.Id),
		PostType:      postTypeStr[post.PostType],
		Status:        postStatusStr[post.Status],
		Title:         post.Title,
		Slug:          post.Slug,
		Content:       post.Content,
		Excerpt:       post.Excerpt,
		CommentStatus: commentStatusStr[post.CommentStatus],
		Locale:        post.Locale,
		HasPassword:   post.PasswordHash != "",
	}
	if post.PublishedAt != nil {
		item.PublishedAt = gtimeToRFC3339(post.PublishedAt)
	}
	if post.CreatedAt != nil {
		item.CreatedAt = gtimeToRFC3339(post.CreatedAt)
	}
	if post.UpdatedAt != nil {
		item.UpdatedAt = gtimeToRFC3339(post.UpdatedAt)
	}

	// 作者
	if post.AuthorId > 0 {
		type AuthorRow struct {
			Id          int64  `orm:"id"`
			Username    string `orm:"username"`
			DisplayName string `orm:"display_name"`
			AvatarId    int64  `orm:"avatar_id"`
		}
		var author AuthorRow
		if e := dao.Users.Ctx(ctx).Where("id", post.AuthorId).Scan(&author); e == nil && author.Id > 0 {
			authorItem := &v1.PostAuthorItem{
				Id:       author.Id,
				Username: author.Username,
				Nickname: author.DisplayName,
			}
			if author.AvatarId > 0 {
				var avatarUrl string
				_ = dao.Medias.Ctx(ctx).Fields("cdn_url").Where("id", author.AvatarId).Scan(&avatarUrl)
				authorItem.Avatar = avatarUrl
			}
			item.Author = authorItem
		}
	}

	// 封面图
	if post.FeaturedImgId != 0 {
		type MediaRow struct {
			Id       int64  `orm:"id"`
			CdnUrl   string `orm:"cdn_url"`
			Filename string `orm:"filename"`
			MimeType string `orm:"mime_type"`
		}
		var media MediaRow
		if e := dao.Medias.Ctx(ctx).Where("id", post.FeaturedImgId).Scan(&media); e == nil && media.Id > 0 {
			item.FeaturedImg = &v1.PostMediaItem{
				Id:       media.Id,
				Url:      media.CdnUrl,
				Title:    media.Filename,
				MimeType: media.MimeType,
			}
		}
	}

	// 统计
	type StatsRow struct {
		ViewCount    int64 `orm:"view_count"`
		CommentCount int64 `orm:"comment_count"`
	}
	var stats StatsRow
	_ = dao.PostStats.Ctx(ctx).Where("post_id", post.Id).Scan(&stats)
	item.ViewCount = stats.ViewCount
	item.CommentCount = stats.CommentCount

	// Record browse history for logged-in users (action = "view", object_type = "post")
	if userID, _ := middleware.GetCurrentUserID(ctx); userID > 0 {
		_, _ = dao.UserActions.Ctx(ctx).Data(g.Map{
			"id":          idgen.New(),
			"user_id":     userID,
			"action":      "view",
			"object_type": "post",
			"object_id":   post.Id,
			"extra":       "{}",
		}).Insert()
	}

	// metas
	item.Metas, _ = s.GetMetas(ctx, int64(post.Id))

	// Run content.render filter — plugins can modify the markdown before it reaches the frontend
	if filtered, err := plugin.Filter(ctx, plugin.FilterContentRender, map[string]any{
		"content": item.Content,
		"type":    item.PostType,
		"id":      item.Id,
		"slug":    item.Slug,
		"title":   item.Title,
	}); err == nil {
		if v, ok := filtered["content"].(string); ok {
			item.Content = v
		}
	}

	return item, nil
}

// ----------------------------------------------------------------
//  GetList  —  分页列表，批量 join
// ----------------------------------------------------------------

func (s *sPost) GetList(ctx context.Context, req *v1.PostGetListReq) (*v1.PostGetListRes, error) {
	// GoFrame auto-adds "WHERE deleted_at IS NULL" to all queries.
	// Trash is represented by status=5; permanently deleted posts have deleted_at set.
	m := dao.Posts.Ctx(ctx)

	if req.Status != nil && *req.Status != "" {
		if si, ok := statusIntMap[*req.Status]; ok {
			m = m.Where("status", si)
		}
	} else {
		// "all" tab: exclude trashed posts (status=5)
		m = m.WhereNot("status", 5)
	}
	if req.PostType != nil && *req.PostType != "" {
		if ti, ok := typeIntMap[*req.PostType]; ok {
			m = m.Where("post_type", ti)
		}
	}
	if req.AuthorId != nil {
		m = m.Where("author_id", *req.AuthorId)
	}
	if req.Keyword != nil && *req.Keyword != "" {
		ids, err := search.Default(ctx).SearchPostIDs(ctx, *req.Keyword)
		if err != nil {
			return nil, err
		}
		if len(ids) == 0 {
			return &v1.PostGetListRes{Data: []*v1.PostListItem{}, Page: req.Page, PageSize: req.PageSize}, nil
		}
		m = m.WhereIn("id", ids)
	}
	if req.TermTaxonomyId != nil {
		m = m.WhereIn("id", dao.ObjectTaxonomies.Ctx(ctx).
			Fields("object_id").
			Where("taxonomy_id", *req.TermTaxonomyId).
			Where("object_type", "post"))
	}
	if req.IncludeCategoryIds != "" {
		ids := parseIdList(req.IncludeCategoryIds)
		if len(ids) > 0 {
			m = m.WhereIn("id", dao.ObjectTaxonomies.Ctx(ctx).
				Fields("object_id").
				WhereIn("taxonomy_id", ids).
				Where("object_type", "post"))
		}
	}
	if req.ExcludeCategoryIds != "" {
		ids := parseIdList(req.ExcludeCategoryIds)
		if len(ids) > 0 {
			m = m.WhereNotIn("id", dao.ObjectTaxonomies.Ctx(ctx).
				Fields("object_id").
				WhereIn("taxonomy_id", ids).
				Where("object_type", "post"))
		}
	}
	if req.MetaKey != nil && *req.MetaKey != "" {
		sub := dao.PostMetas.Ctx(ctx).Fields("post_id").Where("meta_key", *req.MetaKey)
		if req.MetaValue != nil {
			sub = sub.Where("meta_value", *req.MetaValue)
		}
		m = m.WhereIn("id", sub)
	}
	if req.PublishedAfter != nil && *req.PublishedAfter != "" {
		m = m.WhereGTE("published_at", *req.PublishedAfter)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	type PostRow struct {
		Id            int64   `orm:"id"`
		PostType      int     `orm:"post_type"`
		Status        int     `orm:"status"`
		Title         string  `orm:"title"`
		Slug          string  `orm:"slug"`
		Excerpt       string  `orm:"excerpt"`
		CommentStatus int     `orm:"comment_status"`
		Locale        string  `orm:"locale"`
		PublishedAt   *string `orm:"published_at"`
		CreatedAt     string  `orm:"created_at"`
		UpdatedAt     string  `orm:"updated_at"`
		AuthorId      int64   `orm:"author_id"`
		FeaturedImgId *int64  `orm:"featured_img_id"`
	}

	var rows []PostRow
	if total > 0 {
		order := "created_at DESC"
		if req.SortBy != nil && *req.SortBy == "view_count" {
			order = "(SELECT COALESCE(view_count, 0) FROM post_stats WHERE post_id = id) DESC"
		} else if req.SortBy != nil && *req.SortBy == "random" {
			order = "RANDOM()"
		}
		if err = m.Page(req.Page, req.PageSize).Order(order).Scan(&rows); err != nil {
			return nil, err
		}
	}

	// 收集 ID 用于批量加载
	authorIds := make([]int64, 0)
	authorIdSet := make(map[int64]bool)
	imgIds := make([]int64, 0)
	imgIdSet := make(map[int64]bool)
	postIds := make([]int64, 0)

	for _, row := range rows {
		postIds = append(postIds, row.Id)
		if !authorIdSet[row.AuthorId] {
			authorIds = append(authorIds, row.AuthorId)
			authorIdSet[row.AuthorId] = true
		}
		if row.FeaturedImgId != nil && !imgIdSet[*row.FeaturedImgId] {
			imgIds = append(imgIds, *row.FeaturedImgId)
			imgIdSet[*row.FeaturedImgId] = true
		}
	}

	// 批量查作者
	type AuthorRow struct {
		Id          int64  `orm:"id"`
		Username    string `orm:"username"`
		DisplayName string `orm:"display_name"`
	}
	authorMap := make(map[int64]*AuthorRow)
	if len(authorIds) > 0 {
		var authors []AuthorRow
		if err = dao.Users.Ctx(ctx).WhereIn("id", authorIds).Scan(&authors); err != nil {
			return nil, err
		}
		for i := range authors {
			authorMap[authors[i].Id] = &authors[i]
		}
	}

	// 批量查封面图
	type MediaRow struct {
		Id       int64  `orm:"id"`
		CdnUrl   string `orm:"cdn_url"`
		Filename string `orm:"filename"`
		MimeType string `orm:"mime_type"`
	}
	mediaMap := make(map[int64]*MediaRow)
	if len(imgIds) > 0 {
		var medias []MediaRow
		if err = dao.Medias.Ctx(ctx).WhereIn("id", imgIds).Scan(&medias); err != nil {
			return nil, err
		}
		for i := range medias {
			mediaMap[medias[i].Id] = &medias[i]
		}
	}

	// 批量查统计
	type StatsRow struct {
		PostId       int64 `orm:"post_id"`
		ViewCount    int64 `orm:"view_count"`
		CommentCount int64 `orm:"comment_count"`
	}
	statsMap := make(map[int64]*StatsRow)
	if len(postIds) > 0 {
		var stats []StatsRow
		if err = dao.PostStats.Ctx(ctx).
			Fields("post_id, view_count, comment_count").
			WhereIn("post_id", postIds).
			Scan(&stats); err != nil {
			return nil, err
		}
		for i := range stats {
			statsMap[stats[i].PostId] = &stats[i]
		}
	}

	// 批量查 metas
	type MetaRow struct {
		PostId    int64  `orm:"post_id"`
		MetaKey   string `orm:"meta_key"`
		MetaValue string `orm:"meta_value"`
	}
	metasMap := make(map[int64]map[string]string)
	if len(postIds) > 0 {
		var metaRows []MetaRow
		_ = dao.PostMetas.Ctx(ctx).
			Fields("post_id, meta_key, meta_value").
			WhereIn("post_id", postIds).
			Scan(&metaRows)
		for _, mr := range metaRows {
			if metasMap[mr.PostId] == nil {
				metasMap[mr.PostId] = make(map[string]string)
			}
			metasMap[mr.PostId][mr.MetaKey] = mr.MetaValue
		}
	}

	// 组装响应
	list := make([]*v1.PostListItem, len(rows))
	for i, row := range rows {
		item := &v1.PostListItem{
			Id:            row.Id,
			PostType:      postTypeStr[row.PostType],
			Status:        postStatusStr[row.Status],
			Title:         row.Title,
			Slug:          row.Slug,
			Excerpt:       row.Excerpt,
			CommentStatus: commentStatusStr[row.CommentStatus],
			Locale:        row.Locale,
			CreatedAt:     dbStrToRFC3339(row.CreatedAt),
			UpdatedAt:     dbStrToRFC3339(row.UpdatedAt),
		}
		if row.PublishedAt != nil {
			item.PublishedAt = dbStrToRFC3339(*row.PublishedAt)
		}
		if a, ok := authorMap[row.AuthorId]; ok {
			item.Author = &v1.PostAuthorItem{
				Id:       a.Id,
				Username: a.Username,
				Nickname: a.DisplayName,
			}
		}
		if row.FeaturedImgId != nil {
			if med, ok := mediaMap[*row.FeaturedImgId]; ok {
				item.FeaturedImg = &v1.PostMediaItem{
					Id:       med.Id,
					Url:      med.CdnUrl,
					Title:    med.Filename,
					MimeType: med.MimeType,
				}
			}
		}
		if s2, ok := statsMap[row.Id]; ok {
			item.ViewCount = s2.ViewCount
			item.CommentCount = s2.CommentCount
		}
		if m, ok := metasMap[row.Id]; ok {
			item.Metas = m
		}
		list[i] = item
	}

	totalPages := 0
	if req.PageSize > 0 {
		totalPages = (total + req.PageSize - 1) / req.PageSize
	}

	return &v1.PostGetListRes{
		Data:       list,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}
