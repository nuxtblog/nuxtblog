package announcement

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/announcement/v1"
)

type IAnnouncementAdmin interface {
	AnnouncementCreate(ctx context.Context, req *v1.AnnouncementCreateReq) (res *v1.AnnouncementCreateRes, err error)
	AnnouncementListAdmin(ctx context.Context, req *v1.AnnouncementListAdminReq) (res *v1.AnnouncementListAdminRes, err error)
	AnnouncementUpdate(ctx context.Context, req *v1.AnnouncementUpdateReq) (res *v1.AnnouncementUpdateRes, err error)
	AnnouncementDelete(ctx context.Context, req *v1.AnnouncementDeleteReq) (res *v1.AnnouncementDeleteRes, err error)
}

type IAnnouncementPublic interface {
	AnnouncementList(ctx context.Context, req *v1.AnnouncementListReq) (res *v1.AnnouncementListRes, err error)
	AnnouncementMarkRead(ctx context.Context, req *v1.AnnouncementMarkReadReq) (res *v1.AnnouncementMarkReadRes, err error)
}
