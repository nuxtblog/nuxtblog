package media

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/media/v1"
	"github.com/nuxtblog/nuxtblog/internal/consts"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) ExtensionGroupSave(ctx context.Context, req *v1.ExtensionGroupSaveReq) (*v1.ExtensionGroupSaveRes, error) {
	groups := make([]consts.ExtensionGroup, len(req.Groups))
	for i, g := range req.Groups {
		groups[i] = consts.ExtensionGroup{
			Name: g.Name, LabelZh: g.LabelZh, LabelEn: g.LabelEn,
			Extensions: g.Extensions, MaxSizeMB: g.MaxSizeMB,
		}
	}
	if err := service.Media().SaveExtensionGroups(ctx, groups); err != nil {
		return nil, err
	}
	return &v1.ExtensionGroupSaveRes{}, nil
}
