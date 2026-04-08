package post

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/post/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"

	"github.com/gogf/gf/v2/frame/g"
)

func (s *sPost) SeoUpdate(ctx context.Context, req *v1.PostSeoUpdateReq) error {
	_, err := dao.PostSeo.Ctx(ctx).
		Data(g.Map{
			"post_id":         req.Id,
			"meta_title":      req.MetaTitle,
			"meta_desc":       req.MetaDesc,
			"og_title":        req.OgTitle,
			"og_image":        req.OgImage,
			"canonical_url":   req.CanonicalUrl,
			"robots":          req.Robots,
			"structured_data": req.StructuredData,
		}).
		OnConflict("post_id").
		Save()
	return err
}
