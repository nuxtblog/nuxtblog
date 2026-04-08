package doc

import (
	"context"
	"math"

	docv1 "github.com/nuxtblog/nuxtblog/api/doc/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/model/entity"
	plugin "github.com/nuxtblog/nuxtblog/internal/pluginsys"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/nuxtblog/nuxtblog/internal/util/idgen"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sDoc struct{}

func New() *sDoc { return &sDoc{} }

func init() {
	service.RegisterDoc(New())
}

// ----------------------------------------------------------------
//  Collection
// ----------------------------------------------------------------

func (s *sDoc) CollectionCreate(ctx context.Context, req *docv1.CollectionCreateReq) (int64, error) {
	uid, _ := middleware.GetCurrentUserID(ctx)
	col := &entity.DocCollections{
		Slug:        req.Slug,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Locale:      req.Locale,
		SortOrder:   req.SortOrder,
		AuthorId:    int(uid),
	}
	if req.CoverImgId != nil {
		v := int(*req.CoverImgId)
		col.CoverImgId = &v
	}
	return dao.DocCollections.Create(ctx, col)
}

func (s *sDoc) CollectionUpdate(ctx context.Context, req *docv1.CollectionUpdateReq) error {
	data := g.Map{}
	if req.Slug != nil {
		data["slug"] = *req.Slug
	}
	if req.Title != nil {
		data["title"] = *req.Title
	}
	if req.Description != nil {
		data["description"] = *req.Description
	}
	if req.CoverImgId != nil {
		data["cover_img_id"] = *req.CoverImgId
	}
	if req.Status != nil {
		data["status"] = *req.Status
	}
	if req.Locale != nil {
		data["locale"] = *req.Locale
	}
	if req.SortOrder != nil {
		data["sort_order"] = *req.SortOrder
	}
	if len(data) == 0 {
		return nil
	}
	_, err := dao.DocCollections.Ctx(ctx).
		Where("id", req.Id).WhereNull("deleted_at").
		Data(data).Update()
	return err
}

func (s *sDoc) CollectionDelete(ctx context.Context, id int64) error {
	// Check if collection has docs
	count, err := dao.Docs.Ctx(ctx).Where("collection_id", id).WhereNull("deleted_at").Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.NewCode(gcode.CodeInvalidParameter, "collection still has docs, cannot delete")
	}
	_, err = dao.DocCollections.Ctx(ctx).
		Where("id", id).WhereNull("deleted_at").
		Data(g.Map{"deleted_at": gtime.Now()}).Update()
	return err
}

func (s *sDoc) CollectionGetOne(ctx context.Context, id int64) (*docv1.DocCollectionItem, error) {
	type row struct {
		Id          int64       `orm:"id"`
		Slug        string      `orm:"slug"`
		Title       string      `orm:"title"`
		Description string      `orm:"description"`
		CoverImgId  *int64      `orm:"cover_img_id"`
		AuthorId    int64       `orm:"author_id"`
		Status      int         `orm:"status"`
		Locale      string      `orm:"locale"`
		SortOrder   int         `orm:"sort_order"`
		CreatedAt   *gtime.Time `orm:"created_at"`
		UpdatedAt   *gtime.Time `orm:"updated_at"`
	}
	var r row
	if err := dao.DocCollections.Ctx(ctx).Where("id", id).WhereNull("deleted_at").Scan(&r); err != nil {
		return nil, err
	}
	if r.Id == 0 {
		return nil, nil
	}
	docCount, _ := dao.Docs.Ctx(ctx).Where("collection_id", id).WhereNull("deleted_at").Count()
	return &docv1.DocCollectionItem{
		Id:          r.Id,
		Slug:        r.Slug,
		Title:       r.Title,
		Description: r.Description,
		CoverImgId:  r.CoverImgId,
		AuthorId:    r.AuthorId,
		Status:      r.Status,
		Locale:      r.Locale,
		SortOrder:   r.SortOrder,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		DocCount:    docCount,
	}, nil
}

func (s *sDoc) CollectionGetList(ctx context.Context, req *docv1.CollectionGetListReq) (*docv1.CollectionGetListRes, error) {
	m := dao.DocCollections.Ctx(ctx).WhereNull("deleted_at")
	if req.Status != nil {
		m = m.Where("status", *req.Status)
	}
	if req.Locale != nil {
		m = m.Where("locale", *req.Locale)
	}
	if req.AuthorId != nil {
		m = m.Where("author_id", *req.AuthorId)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	type row struct {
		Id          int64       `orm:"id"`
		Slug        string      `orm:"slug"`
		Title       string      `orm:"title"`
		Description string      `orm:"description"`
		CoverImgId  *int64      `orm:"cover_img_id"`
		AuthorId    int64       `orm:"author_id"`
		Status      int         `orm:"status"`
		Locale      string      `orm:"locale"`
		SortOrder   int         `orm:"sort_order"`
		CreatedAt   *gtime.Time `orm:"created_at"`
		UpdatedAt   *gtime.Time `orm:"updated_at"`
	}
	var rows []row
	if total > 0 {
		if err = m.Page(req.Page, req.PageSize).OrderAsc("sort_order").OrderDesc("created_at").Scan(&rows); err != nil {
			return nil, err
		}
	}

	data := make([]*docv1.DocCollectionItem, 0, len(rows))
	for _, r := range rows {
		docCount, _ := dao.Docs.Ctx(ctx).Where("collection_id", r.Id).WhereNull("deleted_at").Count()
		data = append(data, &docv1.DocCollectionItem{
			Id:          r.Id,
			Slug:        r.Slug,
			Title:       r.Title,
			Description: r.Description,
			CoverImgId:  r.CoverImgId,
			AuthorId:    r.AuthorId,
			Status:      r.Status,
			Locale:      r.Locale,
			SortOrder:   r.SortOrder,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
			DocCount:    docCount,
		})
	}

	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	return &docv1.CollectionGetListRes{
		Data:       data,
		Total:      total,
		Page:       req.Page,
		PageSize:   pageSize,
		TotalPages: int(math.Ceil(float64(total) / float64(pageSize))),
	}, nil
}

// ----------------------------------------------------------------
//  Doc
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
			Title    string `orm:"title"`
			Content  string `orm:"content"`
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
		m = m.WhereLike("title", "%"+*req.Keyword+"%")
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

func (s *sDoc) DocSeoUpdate(ctx context.Context, req *docv1.DocSeoUpdateReq) error {
	_, err := dao.DocSeo.Ctx(ctx).Data(g.Map{
		"doc_id":          req.Id,
		"meta_title":      req.MetaTitle,
		"meta_desc":       req.MetaDesc,
		"og_title":        req.OgTitle,
		"og_image":        req.OgImage,
		"canonical_url":   req.CanonicalUrl,
		"robots":          req.Robots,
		"structured_data": req.StructuredData,
	}).Save()
	return err
}

func (s *sDoc) DocRevisionList(ctx context.Context, req *docv1.DocRevisionListReq) (*docv1.DocRevisionListRes, error) {
	m := dao.DocRevisions.Ctx(ctx).Where("doc_id", req.Id)
	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	size := req.Size
	if size < 1 {
		size = 10
	}

	type revRow struct {
		Id        int64       `orm:"id"`
		DocId     int64       `orm:"doc_id"`
		AuthorId  int64       `orm:"author_id"`
		Title     string      `orm:"title"`
		RevNote   string      `orm:"rev_note"`
		CreatedAt *gtime.Time `orm:"created_at"`
	}
	var rows []revRow
	if total > 0 {
		if err = m.Fields("id,doc_id,author_id,title,rev_note,created_at").
			OrderDesc("id").Page(req.Page, size).Scan(&rows); err != nil {
			return nil, err
		}
	}

	list := make([]*docv1.DocRevisionItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, &docv1.DocRevisionItem{
			Id:        r.Id,
			DocId:     r.DocId,
			AuthorId:  r.AuthorId,
			Title:     r.Title,
			RevNote:   r.RevNote,
			CreatedAt: r.CreatedAt,
		})
	}
	return &docv1.DocRevisionListRes{List: list, Total: total}, nil
}

func (s *sDoc) DocRevisionRestore(ctx context.Context, docId, revisionId int64) error {
	type revRow struct {
		Title   string `orm:"title"`
		Content string `orm:"content"`
	}
	var rev revRow
	if err := dao.DocRevisions.Ctx(ctx).Where("id", revisionId).Where("doc_id", docId).Scan(&rev); err != nil {
		return err
	}
	if rev.Title == "" {
		return gerror.NewCode(gcode.CodeNotFound, "revision not found")
	}

	uid, _ := middleware.GetCurrentUserID(ctx)

	// Save current state as a revision before restoring
	type snapRow struct {
		Title   string `orm:"title"`
		Content string `orm:"content"`
	}
	var snap snapRow
	if e := dao.Docs.Ctx(ctx).Where("id", docId).WhereNull("deleted_at").Scan(&snap); e == nil && snap.Title != "" {
		_, _ = dao.DocRevisions.Ctx(ctx).Data(g.Map{
			"id":        idgen.New(),
			"doc_id":    docId,
			"author_id": uid,
			"title":     snap.Title,
			"content":   snap.Content,
			"rev_note":  "pre-restore",
		}).Insert()
	}

	_, err := dao.Docs.Ctx(ctx).Where("id", docId).WhereNull("deleted_at").Data(g.Map{
		"title":   rev.Title,
		"content": rev.Content,
	}).Update()
	return err
}

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
