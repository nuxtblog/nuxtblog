package media

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/media/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) ExtensionGroupList(ctx context.Context, req *v1.ExtensionGroupListReq) (*v1.ExtensionGroupListRes, error) {
	groups, err := service.Media().GetExtensionGroups(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]v1.ExtensionGroupItem, len(groups))
	for i, g := range groups {
		items[i] = v1.ExtensionGroupItem{
			Name: g.Name, LabelZh: g.LabelZh, LabelEn: g.LabelEn,
			Extensions: g.Extensions, MaxSizeMB: g.MaxSizeMB,
		}
	}
	return &v1.ExtensionGroupListRes{List: items}, nil
}
