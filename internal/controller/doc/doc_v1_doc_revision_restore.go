package doc

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/doc/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) DocRevisionRestore(ctx context.Context, req *v1.DocRevisionRestoreReq) (res *v1.DocRevisionRestoreRes, err error) {
	if err = service.Doc().DocRevisionRestore(ctx, req.Id, req.RevisionId); err != nil {
		return nil, err
	}
	return &v1.DocRevisionRestoreRes{}, nil
}
