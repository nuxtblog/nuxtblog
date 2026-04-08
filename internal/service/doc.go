package service

import (
	docv1 "github.com/nuxtblog/nuxtblog/api/doc/v1"
	"context"
)

type IDoc interface {
	CollectionCreate(ctx context.Context, req *docv1.CollectionCreateReq) (int64, error)
	CollectionUpdate(ctx context.Context, req *docv1.CollectionUpdateReq) error
	CollectionDelete(ctx context.Context, id int64) error
	CollectionGetOne(ctx context.Context, id int64) (*docv1.DocCollectionItem, error)
	CollectionGetList(ctx context.Context, req *docv1.CollectionGetListReq) (*docv1.CollectionGetListRes, error)
	DocCreate(ctx context.Context, req *docv1.DocCreateReq) (int64, error)
	DocUpdate(ctx context.Context, req *docv1.DocUpdateReq) error
	DocDelete(ctx context.Context, id int64) error
	DocGetById(ctx context.Context, id int64) (*docv1.DocDetailItem, error)
	DocGetBySlug(ctx context.Context, slug string) (*docv1.DocDetailItem, error)
	DocGetList(ctx context.Context, req *docv1.DocGetListReq) (*docv1.DocGetListRes, error)
	DocSeoUpdate(ctx context.Context, req *docv1.DocSeoUpdateReq) error
	DocRevisionList(ctx context.Context, req *docv1.DocRevisionListReq) (*docv1.DocRevisionListRes, error)
	DocRevisionRestore(ctx context.Context, docId, revisionId int64) error
	IncrementView(ctx context.Context, id int64) error
}

var localDoc IDoc

func Doc() IDoc {
	if localDoc == nil {
		panic("implement not found for interface IDoc, forgot register?")
	}
	return localDoc
}

func RegisterDoc(i IDoc) {
	localDoc = i
}
