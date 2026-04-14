package media

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/media/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) FormatPolicyList(ctx context.Context, req *v1.FormatPolicyListReq) (*v1.FormatPolicyListRes, error) {
	policies, err := service.Media().GetFormatPolicies(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]v1.FormatPolicyItem, len(policies))
	for i, p := range policies {
		items[i] = v1.FormatPolicyItem{
			Name: p.Name, LabelZh: p.LabelZh, LabelEn: p.LabelEn,
			IsSystem: p.IsSystem, Groups: p.Groups,
		}
	}
	return &v1.FormatPolicyListRes{List: items}, nil
}
