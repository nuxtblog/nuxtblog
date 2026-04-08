package option

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/option/v1"
	"github.com/nuxtblog/nuxtblog/internal/storage"
)

func (c *ControllerV1) StorageBackends(ctx context.Context, req *v1.StorageBackendsReq) (*v1.StorageBackendsRes, error) {
	list := storage.ListBackends(ctx)
	backends := make([]v1.StorageBackendInfo, len(list))
	for i, b := range list {
		backends[i] = v1.StorageBackendInfo{
			Name:        b.Name,
			DisplayName: b.DisplayName,
			Type:        b.Type,
			Enabled:     b.Enabled,
		}
	}
	return &v1.StorageBackendsRes{
		Backends: backends,
		Default:  storage.DefaultName(ctx),
	}, nil
}
