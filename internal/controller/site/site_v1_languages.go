package site

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/site/v1"
	"github.com/nuxtblog/nuxtblog/internal/langconf"
)

func (c *ControllerV1) Languages(ctx context.Context, req *v1.LanguagesReq) (res *v1.LanguagesRes, err error) {
	return &v1.LanguagesRes{List: langconf.GetLanguages(ctx)}, nil
}
