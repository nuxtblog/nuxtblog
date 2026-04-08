// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package checkin

import (
	"context"

	"github.com/nuxtblog/nuxtblog/api/checkin/v1"
)

type ICheckinV1 interface {
	DoCheckin(ctx context.Context, req *v1.DoCheckinReq) (res *v1.DoCheckinRes, err error)
	GetCheckinStatus(ctx context.Context, req *v1.GetCheckinStatusReq) (res *v1.GetCheckinStatusRes, err error)
}
