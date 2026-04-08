package verifycode

import (
	"context"

	"github.com/nuxtblog/nuxtblog/api/verifycode/v1"
)

type IVerifycodeV1 interface {
	VerifycodeSend(ctx context.Context, req *v1.VerifycodeSendReq) (res *v1.VerifycodeSendRes, err error)
}
