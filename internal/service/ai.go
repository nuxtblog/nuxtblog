package service

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/ai/v1"
)

type IAI interface {
	ListConfigs(ctx context.Context) (*v1.AIListConfigsRes, error)
	CreateConfig(ctx context.Context, req *v1.AICreateConfigReq) (*v1.AICreateConfigRes, error)
	UpdateConfig(ctx context.Context, req *v1.AIUpdateConfigReq) (*v1.AIUpdateConfigRes, error)
	DeleteConfig(ctx context.Context, id string) error
	ActivateConfig(ctx context.Context, id string) error
	TestConfig(ctx context.Context, id string) (*v1.AITestConfigRes, error)
	// Actions
	Polish(ctx context.Context, content, style string) (string, error)
	Summarize(ctx context.Context, content string, maxLength int) (string, error)
	SuggestTags(ctx context.Context, title, content string) ([]string, error)
	FromURL(ctx context.Context, url, style string) (*v1.AIFromURLRes, error)
	Translate(ctx context.Context, content, targetLang string) (string, error)
}

var _ai IAI

func AI() IAI          { return _ai }
func RegisterAI(a IAI) { _ai = a }
