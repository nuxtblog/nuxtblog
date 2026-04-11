package dao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type friendlinksDao struct{}

var Friendlinks = &friendlinksDao{}

type FriendlinkRow struct {
	Id          int64  `orm:"id"`
	Name        string `orm:"name"`
	Url         string `orm:"url"`
	Logo        string `orm:"logo"`
	Description string `orm:"description"`
	SortOrder   int    `orm:"sort_order"`
	Status      int    `orm:"status"`
	CreatedAt   string `orm:"created_at"`
}

func (d *friendlinksDao) Ctx(ctx context.Context) *gdb.Model {
	return g.DB().Model("friendlinks").Ctx(ctx)
}

func (d *friendlinksDao) Create(ctx context.Context, data g.Map) (int64, error) {
	result, err := d.Ctx(ctx).Insert(data)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (d *friendlinksDao) Update(ctx context.Context, id int64, data g.Map) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	data["updated_at"] = now
	_, err := d.Ctx(ctx).Where("id", id).WhereNull("deleted_at").Update(data)
	return err
}

func (d *friendlinksDao) SoftDelete(ctx context.Context, id int64) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := d.Ctx(ctx).Where("id", id).WhereNull("deleted_at").Update(g.Map{"deleted_at": now})
	return err
}

func (d *friendlinksDao) ListAdmin(ctx context.Context, page, size int) ([]FriendlinkRow, int, error) {
	m := d.Ctx(ctx).WhereNull("deleted_at").OrderAsc("sort_order").OrderDesc("created_at")
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var rows []FriendlinkRow
	if total > 0 {
		err = m.Page(page, size).Scan(&rows)
	}
	return rows, total, err
}

func (d *friendlinksDao) ListPublic(ctx context.Context) ([]FriendlinkRow, error) {
	var rows []FriendlinkRow
	err := d.Ctx(ctx).
		WhereNull("deleted_at").
		Where("status", 1).
		OrderAsc("sort_order").
		OrderDesc("created_at").
		Scan(&rows)
	return rows, err
}
