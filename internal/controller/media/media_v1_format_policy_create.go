package media

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/media/v1"
	"github.com/nuxtblog/nuxtblog/internal/consts"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) FormatPolicyCreate(ctx context.Context, req *v1.FormatPolicyCreateReq) (*v1.FormatPolicyCreateRes, error) {
	if err := service.Media().CreateFormatPolicy(ctx, consts.FormatPolicy{
		Name: req.Name, LabelZh: req.LabelZh, LabelEn: req.LabelEn,
		Groups: req.Groups,
	}); err != nil {
		return nil, err
	}
	return &v1.FormatPolicyCreateRes{}, nil
}
