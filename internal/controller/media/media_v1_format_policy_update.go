package media

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/media/v1"
	"github.com/nuxtblog/nuxtblog/internal/consts"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) FormatPolicyUpdate(ctx context.Context, req *v1.FormatPolicyUpdateReq) (*v1.FormatPolicyUpdateRes, error) {
	if err := service.Media().UpdateFormatPolicy(ctx, req.Name, consts.FormatPolicy{
		LabelZh: req.LabelZh, LabelEn: req.LabelEn,
		Groups: req.Groups,
	}); err != nil {
		return nil, err
	}
	return &v1.FormatPolicyUpdateRes{}, nil
}
