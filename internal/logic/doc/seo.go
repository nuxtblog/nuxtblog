package doc

import (
	"context"

	docv1 "github.com/nuxtblog/nuxtblog/api/doc/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"

	"github.com/gogf/gf/v2/frame/g"
)

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
