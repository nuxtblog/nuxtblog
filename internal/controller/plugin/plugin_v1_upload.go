package plugin

import (
	"context"
	"io"

	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) PluginUploadZip(ctx context.Context, req *v1.PluginUploadZipReq) (res *v1.PluginUploadZipRes, err error) {
	if req.File == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "file is required")
	}
	f, err := req.File.Open()
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "cannot open uploaded file")
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeInternalError, "cannot read file")
	}

	item, err := service.Plugin().InstallZip(ctx, data)
	if err != nil {
		return nil, err
	}
	return &v1.PluginUploadZipRes{Item: *item}, nil
}
