package media

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/media/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) FormatPolicyDelete(ctx context.Context, req *v1.FormatPolicyDeleteReq) (*v1.FormatPolicyDeleteRes, error) {
	if err := service.Media().DeleteFormatPolicy(ctx, req.Name); err != nil {
		return nil, err
	}
	return &v1.FormatPolicyDeleteRes{}, nil
}
