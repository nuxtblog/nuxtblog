package ai

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/ai/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) AIPolish(ctx context.Context, req *v1.AIPolishReq) (*v1.AIPolishRes, error) {
	result, err := service.AI().Polish(ctx, req.Content, req.Style)
	if err != nil {
		return nil, err
	}
	return &v1.AIPolishRes{Result: result}, nil
}

func (c *ControllerV1) AISummarize(ctx context.Context, req *v1.AISummarizeReq) (*v1.AISummarizeRes, error) {
	result, err := service.AI().Summarize(ctx, req.Content, req.MaxLength)
	if err != nil {
		return nil, err
	}
	return &v1.AISummarizeRes{Result: result}, nil
}

func (c *ControllerV1) AISuggestTags(ctx context.Context, req *v1.AISuggestTagsReq) (*v1.AISuggestTagsRes, error) {
	tags, err := service.AI().SuggestTags(ctx, req.Title, req.Content)
	if err != nil {
		return nil, err
	}
	return &v1.AISuggestTagsRes{Tags: tags}, nil
}

func (c *ControllerV1) AIFromURL(ctx context.Context, req *v1.AIFromURLReq) (*v1.AIFromURLRes, error) {
	return service.AI().FromURL(ctx, req.URL, req.Style)
}

func (c *ControllerV1) AITranslate(ctx context.Context, req *v1.AITranslateReq) (*v1.AITranslateRes, error) {
	result, err := service.AI().Translate(ctx, req.Content, req.TargetLang)
	if err != nil {
		return nil, err
	}
	return &v1.AITranslateRes{Result: result}, nil
}
