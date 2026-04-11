package friendlink

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/friendlink/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *AdminControllerV1) FriendlinkCreate(ctx context.Context, req *v1.FriendlinkCreateReq) (*v1.FriendlinkCreateRes, error) {
	id, err := service.Friendlink().Create(ctx, req.Name, req.Url, req.Logo, req.Description, req.SortOrder, req.Status)
	if err != nil {
		return nil, err
	}
	return &v1.FriendlinkCreateRes{Id: id}, nil
}

func (c *AdminControllerV1) FriendlinkAdminList(ctx context.Context, req *v1.FriendlinkAdminListReq) (*v1.FriendlinkAdminListRes, error) {
	list, total, err := service.Friendlink().ListAdmin(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return &v1.FriendlinkAdminListRes{List: list, Total: total}, nil
}

func (c *AdminControllerV1) FriendlinkUpdate(ctx context.Context, req *v1.FriendlinkUpdateReq) (*v1.FriendlinkUpdateRes, error) {
	err := service.Friendlink().Update(ctx, req.Id, req.Name, req.Url, req.Logo, req.Description, req.SortOrder, req.Status)
	if err != nil {
		return nil, err
	}
	return &v1.FriendlinkUpdateRes{}, nil
}

func (c *AdminControllerV1) FriendlinkDelete(ctx context.Context, req *v1.FriendlinkDeleteReq) (*v1.FriendlinkDeleteRes, error) {
	return &v1.FriendlinkDeleteRes{}, service.Friendlink().Delete(ctx, req.Id)
}

func (c *PublicControllerV1) FriendlinkList(ctx context.Context, req *v1.FriendlinkListReq) (*v1.FriendlinkListRes, error) {
	list, err := service.Friendlink().ListPublic(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.FriendlinkListRes{List: list}, nil
}
