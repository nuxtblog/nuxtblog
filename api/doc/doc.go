package doc

import (
	"context"
	"github.com/nuxtblog/nuxtblog/api/doc/v1"
)

type IDocV1 interface {
	CollectionCreate(ctx context.Context, req *v1.CollectionCreateReq) (res *v1.CollectionCreateRes, err error)
	CollectionUpdate(ctx context.Context, req *v1.CollectionUpdateReq) (res *v1.CollectionUpdateRes, err error)
	CollectionDelete(ctx context.Context, req *v1.CollectionDeleteReq) (res *v1.CollectionDeleteRes, err error)
	CollectionGetOne(ctx context.Context, req *v1.CollectionGetOneReq) (res *v1.CollectionGetOneRes, err error)
	CollectionGetList(ctx context.Context, req *v1.CollectionGetListReq) (res *v1.CollectionGetListRes, err error)
	DocCreate(ctx context.Context, req *v1.DocCreateReq) (res *v1.DocCreateRes, err error)
	DocUpdate(ctx context.Context, req *v1.DocUpdateReq) (res *v1.DocUpdateRes, err error)
	DocDelete(ctx context.Context, req *v1.DocDeleteReq) (res *v1.DocDeleteRes, err error)
	DocGetOne(ctx context.Context, req *v1.DocGetOneReq) (res *v1.DocGetOneRes, err error)
	DocGetBySlug(ctx context.Context, req *v1.DocGetBySlugReq) (res *v1.DocGetBySlugRes, err error)
	DocGetList(ctx context.Context, req *v1.DocGetListReq) (res *v1.DocGetListRes, err error)
	DocSeoUpdate(ctx context.Context, req *v1.DocSeoUpdateReq) (res *v1.DocSeoUpdateRes, err error)
	DocRevisionList(ctx context.Context, req *v1.DocRevisionListReq) (res *v1.DocRevisionListRes, err error)
	DocRevisionRestore(ctx context.Context, req *v1.DocRevisionRestoreReq) (res *v1.DocRevisionRestoreRes, err error)
	DocView(ctx context.Context, req *v1.DocViewReq) (res *v1.DocViewRes, err error)
}
