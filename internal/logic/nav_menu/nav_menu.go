package nav_menu

import (
	"context"
	"errors"
	"strings"

	v1 "github.com/nuxtblog/nuxtblog/api/nav_menu/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/model/entity"
	"github.com/nuxtblog/nuxtblog/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sNavMenu struct{}

func New() service.INavMenu { return &sNavMenu{} }

func init() {
	service.RegisterNavMenu(New())
}

func (s *sNavMenu) List(ctx context.Context) ([]*v1.NavMenuOutput, error) {
	var menus []entity.NavMenus
	if err := dao.NavMenus.Ctx(ctx).OrderAsc("id").Scan(&menus); err != nil {
		return nil, err
	}
	if len(menus) == 0 {
		return []*v1.NavMenuOutput{}, nil
	}

	menuIDs := make([]int, 0, len(menus))
	for _, m := range menus {
		menuIDs = append(menuIDs, m.Id)
	}

	var items []entity.NavMenuItems
	cols := dao.NavMenuItems.Columns()
	if err := dao.NavMenuItems.Ctx(ctx).
		WhereIn(cols.MenuId, menuIDs).
		OrderAsc(cols.SortOrder).
		Scan(&items); err != nil && !strings.Contains(err.Error(), "no rows") {
		return nil, err
	}

	// Index items by menu_id
	itemsByMenu := make(map[int][]*v1.NavMenuItemOutput)
	for i := range items {
		it := &items[i]
		out := itemToOutput(it)
		itemsByMenu[it.MenuId] = append(itemsByMenu[it.MenuId], out)
	}

	out := make([]*v1.NavMenuOutput, 0, len(menus))
	for i := range menus {
		m := &menus[i]
		its := itemsByMenu[m.Id]
		if its == nil {
			its = []*v1.NavMenuItemOutput{}
		}
		out = append(out, menuToOutput(m, its))
	}
	return out, nil
}

func (s *sNavMenu) Create(ctx context.Context, req *v1.NavMenuCreateReq) (*v1.NavMenuOutput, error) {
	now := gtime.Now()
	id, err := dao.NavMenus.Ctx(ctx).InsertAndGetId(g.Map{
		"name":        req.Name,
		"location":    req.Location,
		"description": req.Description,
		"created_at":  now,
		"updated_at":  now,
	})
	if err != nil {
		return nil, err
	}
	return s.GetOne(ctx, id)
}

func (s *sNavMenu) GetOne(ctx context.Context, id int64) (*v1.NavMenuOutput, error) {
	var m entity.NavMenus
	if err := dao.NavMenus.Ctx(ctx).Where("id", id).Scan(&m); err != nil {
		return nil, err
	}
	if m.Id == 0 {
		return nil, errors.New(g.I18n().T(ctx, "error.menu_not_found"))
	}

	cols := dao.NavMenuItems.Columns()
	var items []entity.NavMenuItems
	if err := dao.NavMenuItems.Ctx(ctx).
		Where(cols.MenuId, id).
		OrderAsc(cols.SortOrder).
		Scan(&items); err != nil && !strings.Contains(err.Error(), "no rows") {
		return nil, err
	}

	outItems := make([]*v1.NavMenuItemOutput, 0, len(items))
	for i := range items {
		outItems = append(outItems, itemToOutput(&items[i]))
	}
	return menuToOutput(&m, outItems), nil
}

func (s *sNavMenu) Update(ctx context.Context, req *v1.NavMenuUpdateReq) (*v1.NavMenuOutput, error) {
	// 1. Update menu meta
	data := g.Map{"updated_at": gtime.Now()}
	if req.Name != nil {
		data["name"] = *req.Name
	}
	if req.Location != nil {
		data["location"] = *req.Location
	}
	if req.Description != nil {
		data["description"] = *req.Description
	}
	if _, err := dao.NavMenus.Ctx(ctx).Where("id", req.Id).Data(data).Update(); err != nil {
		return nil, err
	}

	// 2. Replace items
	cols := dao.NavMenuItems.Columns()
	if _, err := dao.NavMenuItems.Ctx(ctx).Where(cols.MenuId, req.Id).Delete(); err != nil {
		return nil, err
	}

	if len(req.Items) > 0 {
		// First pass: insert all items (parent_id=0), collect index->id mapping
		insertedIDs := make([]int64, len(req.Items))
		for i, item := range req.Items {
			id, err := dao.NavMenuItems.Ctx(ctx).InsertAndGetId(g.Map{
				cols.MenuId:     req.Id,
				cols.ParentId:   0,
				cols.ObjectType: item.ObjectType,
				cols.ObjectId:   item.ObjectId,
				cols.Label:      item.Label,
				cols.Url:        item.Url,
				cols.Target:     item.Target,
				cols.CssClasses: item.CssClasses,
				cols.SortOrder:  i,
			})
			if err != nil {
				return nil, err
			}
			insertedIDs[i] = id
		}

		// Second pass: update parent_id for items that have a parent
		for i, item := range req.Items {
			if item.ParentIdx >= 0 && item.ParentIdx < len(insertedIDs) {
				parentID := insertedIDs[item.ParentIdx]
				if _, err := dao.NavMenuItems.Ctx(ctx).
					Where(cols.Id, insertedIDs[i]).
					Data(g.Map{cols.ParentId: parentID}).
					Update(); err != nil {
					return nil, err
				}
			}
		}
	}

	return s.GetOne(ctx, req.Id)
}

func (s *sNavMenu) Delete(ctx context.Context, id int64) error {
	_, err := dao.NavMenus.Ctx(ctx).Where("id", id).Delete()
	return err
}

func (s *sNavMenu) GetByLocation(ctx context.Context, location string) (*v1.NavMenuOutput, error) {
	cols := dao.NavMenus.Columns()
	var m entity.NavMenus
	if err := dao.NavMenus.Ctx(ctx).Where(cols.Location, location).Scan(&m); err != nil {
		return nil, err
	}
	if m.Id == 0 {
		return nil, nil
	}
	return s.GetOne(ctx, int64(m.Id))
}

func menuToOutput(m *entity.NavMenus, items []*v1.NavMenuItemOutput) *v1.NavMenuOutput {
	return &v1.NavMenuOutput{
		Id:          int64(m.Id),
		Name:        m.Name,
		Location:    m.Location,
		Description: m.Description,
		Items:       items,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func itemToOutput(it *entity.NavMenuItems) *v1.NavMenuItemOutput {
	return &v1.NavMenuItemOutput{
		Id:         int64(it.Id),
		MenuId:     int64(it.MenuId),
		ParentId:   int64(it.ParentId),
		ObjectType: it.ObjectType,
		ObjectId:   int64(it.ObjectId),
		Label:      it.Label,
		Url:        it.Url,
		Target:     it.Target,
		CssClasses: it.CssClasses,
		SortOrder:  it.SortOrder,
	}
}
