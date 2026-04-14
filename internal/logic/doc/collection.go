package doc

import (
	"context"
	"math"

	docv1 "github.com/nuxtblog/nuxtblog/api/doc/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/model/entity"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

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
