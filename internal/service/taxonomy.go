package service

import (
	v1 "github.com/nuxtblog/nuxtblog/api/taxonomy/v1"
	"context"
)

type (
	ITaxonomy interface {
		TaxonomyCreate(ctx context.Context, req *v1.TaxonomyCreateReq) (int64, error)
		TaxonomyDelete(ctx context.Context, id int64) error
		TaxonomyGetList(ctx context.Context, req *v1.TaxonomyGetListReq) (*v1.TaxonomyGetListRes, error)
		TaxonomyGetTree(ctx context.Context, taxonomy string) (*v1.TaxonomyGetTreeRes, error)
		TaxonomyUpdate(ctx context.Context, req *v1.TaxonomyUpdateReq) error
		TermCreate(ctx context.Context, req *v1.TermCreateReq) (int64, error)
		TermDelete(ctx context.Context, id int64) error
		TermUpdate(ctx context.Context, req *v1.TermUpdateReq) error
		ObjectTaxonomyBind(ctx context.Context, req *v1.ObjectTaxonomyBindReq) error
		ObjectTaxonomyUnbind(ctx context.Context, req *v1.ObjectTaxonomyUnbindReq) error
		ObjectTaxonomyGet(ctx context.Context, req *v1.ObjectTaxonomyGetReq) (*v1.ObjectTaxonomyGetRes, error)
	}
)

var (
	localTaxonomy ITaxonomy
)

func Taxonomy() ITaxonomy {
	if localTaxonomy == nil {
		panic("implement not found for interface ITaxonomy, forgot register?")
	}
	return localTaxonomy
}

func RegisterTaxonomy(i ITaxonomy) {
	localTaxonomy = i
}
