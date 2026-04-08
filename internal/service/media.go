package service

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/media/v1"
)

type (
	IMedia interface {
		Upload(ctx context.Context, req *v1.MediaUploadReq) (*v1.MediaItem, error)
		Link(ctx context.Context, req *v1.MediaLinkReq) (*v1.MediaItem, error)
		Delete(ctx context.Context, id int64) error
		Update(ctx context.Context, req *v1.MediaUpdateReq) error
		GetOne(ctx context.Context, id int64) (*v1.MediaItem, error)
		GetList(ctx context.Context, req *v1.MediaGetListReq) (*v1.MediaGetListRes, error)
		GetStats(ctx context.Context) (*v1.MediaGetStatsRes, error)
		Localize(ctx context.Context, req *v1.MediaLocalizeReq) (*v1.MediaItem, error)
		// Category
		GetCategories(ctx context.Context) ([]v1.MediaCategoryItem, error)
		UpdateCategory(ctx context.Context, req *v1.MediaCategoryUpdateReq) error
	}
)

var (
	localMedia IMedia
)

func Media() IMedia {
	if localMedia == nil {
		panic("implement not found for interface IMedia, forgot register?")
	}
	return localMedia
}

func RegisterMedia(i IMedia) {
	localMedia = i
}
