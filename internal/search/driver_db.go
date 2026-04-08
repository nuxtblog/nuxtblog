package search

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
)

// dbSearcher uses SQL LIKE on the posts table.
// Works with both SQLite and PostgreSQL — no extra setup required.
type dbSearcher struct{}

func newDbSearcher() Searcher { return &dbSearcher{} }

func (s *dbSearcher) SearchPostIDs(ctx context.Context, keyword string) ([]int64, error) {
	type Row struct {
		Id int64 `orm:"id"`
	}
	var rows []Row
	err := g.DB().Ctx(ctx).
		Model("posts").
		Fields("id").
		WhereNull("deleted_at").
		WhereLike("title", fmt.Sprintf("%%%s%%", keyword)).
		Limit(1000).
		Scan(&rows)
	if err != nil {
		return nil, err
	}
	ids := make([]int64, len(rows))
	for i, r := range rows {
		ids[i] = r.Id
	}
	return ids, nil
}

// IndexPost is a no-op: the database IS the index.
func (s *dbSearcher) IndexPost(_ context.Context, _ PostDoc) error { return nil }

// DeletePost is a no-op: deletion from the DB removes the post automatically.
func (s *dbSearcher) DeletePost(_ context.Context, _ int64) error { return nil }
