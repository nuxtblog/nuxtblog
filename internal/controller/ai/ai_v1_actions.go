package ai

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/ai/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) AIFromURL(ctx context.Context, req *v1.AIFromURLReq) (*v1.AIFromURLRes, error) {
	return service.AI().FromURL(ctx, req.URL, req.Style)
}
