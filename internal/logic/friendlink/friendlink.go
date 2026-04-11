package friendlink

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/friendlink/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/gogf/gf/v2/frame/g"
)

type sFriendlink struct{}

func New() service.IFriendlink { return &sFriendlink{} }

func init() {
	service.RegisterFriendlink(New())
}

func (s *sFriendlink) Create(ctx context.Context, name, url, logo, description string, sortOrder, status int) (int64, error) {
	return dao.Friendlinks.Create(ctx, g.Map{
		"name":        name,
		"url":         url,
		"logo":        logo,
		"description": description,
		"sort_order":  sortOrder,
		"status":      status,
	})
}

func (s *sFriendlink) ListAdmin(ctx context.Context, page, size int) ([]*v1.FriendlinkItem, int, error) {
	rows, total, err := dao.Friendlinks.ListAdmin(ctx, page, size)
	if err != nil {
		return nil, 0, err
	}
	list := make([]*v1.FriendlinkItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, rowToItem(row))
	}
	return list, total, nil
}

func (s *sFriendlink) Update(ctx context.Context, id int64, name, url, logo, description string, sortOrder, status int) error {
	return dao.Friendlinks.Update(ctx, id, g.Map{
		"name":        name,
		"url":         url,
		"logo":        logo,
		"description": description,
		"sort_order":  sortOrder,
		"status":      status,
	})
}

func (s *sFriendlink) Delete(ctx context.Context, id int64) error {
	return dao.Friendlinks.SoftDelete(ctx, id)
}

func (s *sFriendlink) ListPublic(ctx context.Context) ([]*v1.FriendlinkItem, error) {
	rows, err := dao.Friendlinks.ListPublic(ctx)
	if err != nil {
		return nil, err
	}
	list := make([]*v1.FriendlinkItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, rowToItem(row))
	}
	return list, nil
}

func rowToItem(row dao.FriendlinkRow) *v1.FriendlinkItem {
	return &v1.FriendlinkItem{
		Id:          row.Id,
		Name:        row.Name,
		Url:         row.Url,
		Logo:        row.Logo,
		Description: row.Description,
		SortOrder:   row.SortOrder,
		Status:      row.Status,
		CreatedAt:   row.CreatedAt,
	}
}
