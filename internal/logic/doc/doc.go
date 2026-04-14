package doc

import (
	"context"
	"math"

	docv1 "github.com/nuxtblog/nuxtblog/api/doc/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/model/entity"
	plugin "github.com/nuxtblog/nuxtblog/internal/pluginsys"
	"github.com/nuxtblog/nuxtblog/internal/search"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/nuxtblog/nuxtblog/internal/util/idgen"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sDoc struct{}

func New() *sDoc { return &sDoc{} }

func init() {
	service.RegisterDoc(New())
}

// ----------------------------------------------------------------
//  Doc CRUD
// ----------------------------------------------------------------

func (s *sDoc) DocCreate(ctx context.Context, req *docv1.DocCreateReq) (int64, error) {
	uid, _ := middleware.GetCurrentUserID(ctx)
	doc := &entity.Docs{
		CollectionId:  int(req.CollectionId),
		Title:         req.Title,
		Slug:          req.Slug,
		Content:       req.Content,
		Excerpt:       req.Excerpt,
		Status:        req.Status,
		CommentStatus: req.CommentStatus,
		Locale:        req.Locale,
		SortOrder:     req.SortOrder,
		AuthorId:      int(uid),
		PublishedAt:   req.PublishedAt,
	}
	if req.ParentId != nil {
		v := int(*req.ParentId)
		doc.ParentId = &v
	}

	id, err := dao.Docs.Create(ctx, doc)
	if err != nil {
		return 0, err
	}

	// Save initial revision
	_, _ = dao.DocRevisions.Ctx(ctx).Data(g.Map{
		"doc_id":    id,
		"author_id": uid,
		"title":     req.Title,
		"content":   req.Content,
		"rev_note":  "initial",
	}).Insert()

	// Insert initial doc_stats
	_, _ = dao.DocStats.Ctx(ctx).Data(g.Map{
		"doc_id":        id,
		"view_count":    0,
		"like_count":    0,
		"comment_count": 0,
	}).Insert()

	return id, nil
}

func (s *sDoc) DocUpdate(ctx context.Context, req *docv1.DocUpdateReq) error {
	data := g.Map{}
	if req.CollectionId != nil {
		data["collection_id"] = *req.CollectionId
	}
	if req.ParentId != nil {
		if *req.ParentId == 0 {
			data["parent_id"] = nil
		} else {
			data["parent_id"] = *req.ParentId
		}
	}
	if req.Title != nil {
		data["title"] = *req.Title
	}
	if req.Slug != nil {
		data["slug"] = *req.Slug
	}
	if req.Content != nil {
		data["content"] = *req.Content
	}
	if req.Excerpt != nil {
		data["excerpt"] = *req.Excerpt
	}
	if req.Status != nil {
		data["status"] = *req.Status
	}
	if req.CommentStatus != nil {
		data["comment_status"] = *req.CommentStatus
	}
	if req.Locale != nil {
		data["locale"] = *req.Locale
	}
	if req.SortOrder != nil {
		data["sort_order"] = *req.SortOrder
	}
	if req.PublishedAt != nil {
		data["published_at"] = req.PublishedAt
	}
	if len(data) == 0 {
		return nil
	}

	// Save revision if title or content changed
	if req.Title != nil || req.Content != nil {
		uid, _ := middleware.GetCurrentUserID(ctx)
		type snapRow struct {
			Title   string `orm:"title"`
			Content string `orm:"content"`
		}
		var snap snapRow
		if e := dao.Docs.Ctx(ctx).Where("id", req.Id).WhereNull("deleted_at").Scan(&snap); e == nil && snap.Title != "" {
			_, _ = dao.DocRevisions.Ctx(ctx).Data(g.Map{
				"id":        idgen.New(),
				"doc_id":    req.Id,
				"author_id": uid,
				"title":     snap.Title,
				"content":   snap.Content,
				"rev_note":  "auto",
			}).Insert()
		}
	}

	_, err := dao.Docs.Ctx(ctx).Where("id", req.Id).WhereNull("deleted_at").Data(data).Update()
	return err
}

func (s *sDoc) DocDelete(ctx context.Context, id int64) error {
	_, err := dao.Docs.Ctx(ctx).
		Where("id", id).WhereNull("deleted_at").
		Data(g.Map{"deleted_at": gtime.Now()}).Update()
	return err
}

// ----------------------------------------------------------------
//  Doc Query
// ----------------------------------------------------------------

func (s *sDoc) docScanToDetail(ctx context.Context, docId int64, d *entity.Docs) (*docv1.DocDetailItem, error) {
	item := &docv1.DocDetailItem{
		DocItem: docv1.DocItem{
			Id:            int64(d.Id),
			CollectionId:  int64(d.CollectionId),
			SortOrder:     d.SortOrder,
			Status:        d.Status,
			Title:         d.Title,
			Slug:          d.Slug,
			Excerpt:       d.Excerpt,
			AuthorId:      int64(d.AuthorId),
			CommentStatus: d.CommentStatus,
			Locale:        d.Locale,
			PublishedAt:   d.PublishedAt,
			CreatedAt:     d.CreatedAt,
			UpdatedAt:     d.UpdatedAt,
		},
		Content: d.Content,
	}
	if d.ParentId != nil {
		v := int64(*d.ParentId)
		item.DocItem.ParentId = &v
	}

	// Load stats
	type statsRow struct {
		ViewCount    int64 `orm:"view_count"`
		LikeCount    int64 `orm:"like_count"`
		CommentCount int64 `orm:"comment_count"`
	}
	var stats statsRow
	if err := dao.DocStats.Ctx(ctx).Where("doc_id", docId).Scan(&stats); err == nil {
		item.DocItem.Stats = &docv1.DocStatsItem{
			ViewCount:    stats.ViewCount,
			LikeCount:    stats.LikeCount,
			CommentCount: stats.CommentCount,
		}
	}

	// Load SEO
	type seoRow struct {
		MetaTitle      string `orm:"meta_title"`
		MetaDesc       string `orm:"meta_desc"`
		OgTitle        string `orm:"og_title"`
		OgImage        string `orm:"og_image"`
		CanonicalUrl   string `orm:"canonical_url"`
		Robots         string `orm:"robots"`
		StructuredData string `orm:"structured_data"`
	}
	var seo seoRow
	if err := dao.DocSeo.Ctx(ctx).Where("doc_id", docId).Scan(&seo); err == nil && seo.MetaTitle != "" {
		item.Seo = &docv1.DocSeoItem{
			MetaTitle:      seo.MetaTitle,
			MetaDesc:       seo.MetaDesc,
			OgTitle:        seo.OgTitle,
			OgImage:        seo.OgImage,
			CanonicalUrl:   seo.CanonicalUrl,
			Robots:         seo.Robots,
			StructuredData: seo.StructuredData,
		}
	}

	// Run content.render filter — plugins can modify the markdown before it reaches the frontend
	if filtered, err := plugin.Filter(ctx, plugin.FilterContentRender, map[string]any{
		"content": item.Content,
		"type":    "doc",
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

func (s *sDoc) DocGetById(ctx context.Context, id int64) (*docv1.DocDetailItem, error) {
	d, err := dao.Docs.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, nil
	}
	return s.docScanToDetail(ctx, id, d)
}

func (s *sDoc) DocGetBySlug(ctx context.Context, slug string) (*docv1.DocDetailItem, error) {
	d, err := dao.Docs.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, nil
	}
	return s.docScanToDetail(ctx, int64(d.Id), d)
}

func (s *sDoc) DocGetList(ctx context.Context, req *docv1.DocGetListReq) (*docv1.DocGetListRes, error) {
	m := dao.Docs.Ctx(ctx).WhereNull("deleted_at")
	if req.CollectionId != nil {
		m = m.Where("collection_id", *req.CollectionId)
	}
	if req.ParentId != nil {
		if *req.ParentId == 0 {
			m = m.WhereNull("parent_id")
		} else {
			m = m.Where("parent_id", *req.ParentId)
		}
	}
	if req.Status != nil {
		m = m.Where("status", *req.Status)
	}
	if req.AuthorId != nil {
		m = m.Where("author_id", *req.AuthorId)
	}
	if req.Locale != nil {
		m = m.Where("locale", *req.Locale)
	}
	if req.Keyword != nil && *req.Keyword != "" {
		sr, err := search.Default(ctx).Search(ctx, search.ContentDoc, *req.Keyword)
		if err != nil {
			return nil, err
		}
		if len(sr.IDs) == 0 {
			return &docv1.DocGetListRes{Data: []*docv1.DocItem{}, Page: req.Page, PageSize: req.PageSize}, nil
		}
		m = m.WhereIn("id", sr.IDs)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	type docRow struct {
		Id            int64       `orm:"id"`
		CollectionId  int64       `orm:"collection_id"`
		ParentId      *int64      `orm:"parent_id"`
		SortOrder     int         `orm:"sort_order"`
		Status        int         `orm:"status"`
		Title         string      `orm:"title"`
		Slug          string      `orm:"slug"`
		Excerpt       string      `orm:"excerpt"`
		AuthorId      int64       `orm:"author_id"`
		CommentStatus int         `orm:"comment_status"`
		Locale        string      `orm:"locale"`
		PublishedAt   *gtime.Time `orm:"published_at"`
		CreatedAt     *gtime.Time `orm:"created_at"`
		UpdatedAt     *gtime.Time `orm:"updated_at"`
	}
	var rows []docRow
	if total > 0 {
		if err = m.Page(req.Page, req.PageSize).OrderAsc("sort_order").OrderDesc("created_at").Scan(&rows); err != nil {
			return nil, err
		}
	}

	data := make([]*docv1.DocItem, 0, len(rows))
	for _, r := range rows {
		item := &docv1.DocItem{
			Id:            r.Id,
			CollectionId:  r.CollectionId,
			ParentId:      r.ParentId,
			SortOrder:     r.SortOrder,
			Status:        r.Status,
			Title:         r.Title,
			Slug:          r.Slug,
			Excerpt:       r.Excerpt,
			AuthorId:      r.AuthorId,
			CommentStatus: r.CommentStatus,
			Locale:        r.Locale,
			PublishedAt:   r.PublishedAt,
			CreatedAt:     r.CreatedAt,
			UpdatedAt:     r.UpdatedAt,
		}
		// Load stats
		type statsRow struct {
			ViewCount    int64 `orm:"view_count"`
			LikeCount    int64 `orm:"like_count"`
			CommentCount int64 `orm:"comment_count"`
		}
		var stats statsRow
		if e := dao.DocStats.Ctx(ctx).Where("doc_id", r.Id).Scan(&stats); e == nil {
			item.Stats = &docv1.DocStatsItem{
				ViewCount:    stats.ViewCount,
				LikeCount:    stats.LikeCount,
				CommentCount: stats.CommentCount,
			}
		}
		data = append(data, item)
	}

	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	return &docv1.DocGetListRes{
		Data:       data,
		Total:      total,
		Page:       req.Page,
		PageSize:   pageSize,
		TotalPages: int(math.Ceil(float64(total) / float64(pageSize))),
	}, nil
}

// ----------------------------------------------------------------
//  View
// ----------------------------------------------------------------

func (s *sDoc) IncrementView(ctx context.Context, id int64) error {
	_, err := dao.DocStats.DB().Exec(
		ctx,
		`INSERT INTO doc_stats (doc_id, view_count, like_count, comment_count)
		 VALUES (?, 1, 0, 0)
		 ON CONFLICT(doc_id) DO UPDATE SET view_count = view_count + 1`,
		id,
	)
	if err != nil {
		g.Log().Warningf(ctx, "[doc] increment view error: %v", err)
	}
	return nil
}
