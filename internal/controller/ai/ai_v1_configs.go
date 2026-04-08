package ai

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/ai/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) AIListConfigs(ctx context.Context, req *v1.AIListConfigsReq) (*v1.AIListConfigsRes, error) {
	return service.AI().ListConfigs(ctx)
}

func (c *ControllerV1) AICreateConfig(ctx context.Context, req *v1.AICreateConfigReq) (*v1.AICreateConfigRes, error) {
	return service.AI().CreateConfig(ctx, req)
}

func (c *ControllerV1) AIUpdateConfig(ctx context.Context, req *v1.AIUpdateConfigReq) (*v1.AIUpdateConfigRes, error) {
	return service.AI().UpdateConfig(ctx, req)
}

func (c *ControllerV1) AIDeleteConfig(ctx context.Context, req *v1.AIDeleteConfigReq) (*v1.AIDeleteConfigRes, error) {
	return &v1.AIDeleteConfigRes{}, service.AI().DeleteConfig(ctx, req.ID)
}

func (c *ControllerV1) AIActivateConfig(ctx context.Context, req *v1.AIActivateConfigReq) (*v1.AIActivateConfigRes, error) {
	return &v1.AIActivateConfigRes{}, service.AI().ActivateConfig(ctx, req.ID)
}

func (c *ControllerV1) AITestConfig(ctx context.Context, req *v1.AITestConfigReq) (*v1.AITestConfigRes, error) {
	return service.AI().TestConfig(ctx, req.ID)
}
