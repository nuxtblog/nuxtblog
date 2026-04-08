// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package taxonomy

import (
	"context"

	"github.com/nuxtblog/nuxtblog/api/taxonomy/v1"
)

type ITaxonomyV1 interface {
	TermCreate(ctx context.Context, req *v1.TermCreateReq) (res *v1.TermCreateRes, err error)
	TermDelete(ctx context.Context, req *v1.TermDeleteReq) (res *v1.TermDeleteRes, err error)
	TermUpdate(ctx context.Context, req *v1.TermUpdateReq) (res *v1.TermUpdateRes, err error)
	TaxonomyCreate(ctx context.Context, req *v1.TaxonomyCreateReq) (res *v1.TaxonomyCreateRes, err error)
	TaxonomyDelete(ctx context.Context, req *v1.TaxonomyDeleteReq) (res *v1.TaxonomyDeleteRes, err error)
	TaxonomyUpdate(ctx context.Context, req *v1.TaxonomyUpdateReq) (res *v1.TaxonomyUpdateRes, err error)
	TaxonomyGetTree(ctx context.Context, req *v1.TaxonomyGetTreeReq) (res *v1.TaxonomyGetTreeRes, err error)
	TaxonomyGetList(ctx context.Context, req *v1.TaxonomyGetListReq) (res *v1.TaxonomyGetListRes, err error)
	ObjectTaxonomyBind(ctx context.Context, req *v1.ObjectTaxonomyBindReq) (res *v1.ObjectTaxonomyBindRes, err error)
	ObjectTaxonomyUnbind(ctx context.Context, req *v1.ObjectTaxonomyUnbindReq) (res *v1.ObjectTaxonomyUnbindRes, err error)
	ObjectTaxonomyGet(ctx context.Context, req *v1.ObjectTaxonomyGetReq) (res *v1.ObjectTaxonomyGetRes, err error)
}
