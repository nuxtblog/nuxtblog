package ai

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/ai/v1"
)

type IAI interface {
	AIListConfigs(ctx context.Context, req *v1.AIListConfigsReq) (res *v1.AIListConfigsRes, err error)
	AICreateConfig(ctx context.Context, req *v1.AICreateConfigReq) (res *v1.AICreateConfigRes, err error)
	AIUpdateConfig(ctx context.Context, req *v1.AIUpdateConfigReq) (res *v1.AIUpdateConfigRes, err error)
	AIDeleteConfig(ctx context.Context, req *v1.AIDeleteConfigReq) (res *v1.AIDeleteConfigRes, err error)
	AIActivateConfig(ctx context.Context, req *v1.AIActivateConfigReq) (res *v1.AIActivateConfigRes, err error)
	AITestConfig(ctx context.Context, req *v1.AITestConfigReq) (res *v1.AITestConfigRes, err error)
	AIPolish(ctx context.Context, req *v1.AIPolishReq) (res *v1.AIPolishRes, err error)
	AISummarize(ctx context.Context, req *v1.AISummarizeReq) (res *v1.AISummarizeRes, err error)
	AISuggestTags(ctx context.Context, req *v1.AISuggestTagsReq) (res *v1.AISuggestTagsRes, err error)
	AIFromURL(ctx context.Context, req *v1.AIFromURLReq) (res *v1.AIFromURLRes, err error)
	AITranslate(ctx context.Context, req *v1.AITranslateReq) (res *v1.AITranslateRes, err error)
}
