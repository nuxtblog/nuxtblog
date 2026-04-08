package site

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/site/v1"
)

type ISiteV1 interface {
	Languages(ctx context.Context, req *v1.LanguagesReq) (res *v1.LanguagesRes, err error)
}
