package verifycode

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/verifycode/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) VerifycodeSend(ctx context.Context, req *v1.VerifycodeSendReq) (res *v1.VerifycodeSendRes, err error) {
	return service.Verifycode().Send(ctx, req)
}
