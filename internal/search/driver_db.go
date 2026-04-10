package search

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
)

// dbSearcher uses SQL LIKE on the relevant table.
// Works with both SQLite and PostgreSQL — no extra setup required.
type dbSearcher struct{}

func newDbSearcher() Searcher { return &dbSearcher{} }

func (s *dbSearcher) Search(ctx context.Context, ct ContentType, keyword string) (SearchResult, error) {
	meta := metaMap[ct]
	type Row struct {
		Id int64 `orm:"id"`
	}
	var rows []Row
	err := g.DB().Ctx(ctx).
		Model(meta.Table).
		Fields("id").
		WhereNull("deleted_at").
		WhereLike(meta.LikeCol, fmt.Sprintf("%%%s%%", keyword)).
		Limit(defaultSearchLimit).
		Scan(&rows)
	if err != nil {
		return SearchResult{}, err
	}
	ids := make([]int64, len(rows))
	for i, r := range rows {
		ids[i] = r.Id
	}
	return SearchResult{IDs: ids, Total: -1}, nil
}

// Index is a no-op: the database IS the index.
func (s *dbSearcher) Index(_ context.Context, _ ContentType, _ Document) error { return nil }

// Delete is a no-op: deletion from the DB removes the record automatically.
func (s *dbSearcher) Delete(_ context.Context, _ ContentType, _ int64) error { return nil }
