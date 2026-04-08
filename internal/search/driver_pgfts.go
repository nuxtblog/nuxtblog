package search

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

// pgFtsSearcher uses PostgreSQL full-text search (tsvector / to_tsquery).
// Requires PostgreSQL — will NOT work with SQLite.
// Supports Chinese via pg_jieba or zhparser extensions (optional).
type pgFtsSearcher struct{}

func newPgFtsSearcher() Searcher { return &pgFtsSearcher{} }

func (s *pgFtsSearcher) SearchPostIDs(ctx context.Context, keyword string) ([]int64, error) {
	type Row struct {
		Id int64 `orm:"id"`
	}
	var rows []Row
	// plainto_tsquery handles multi-word input safely (no need to escape).
	// Uses the "simple" dictionary so it works without language-specific configs.
	err := g.DB().Ctx(ctx).
		Model("posts").
		Fields("id").
		WhereNull("deleted_at").
		Where("to_tsvector('simple', title || ' ' || excerpt) @@ plainto_tsquery('simple', ?)", keyword).
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

func (s *pgFtsSearcher) IndexPost(_ context.Context, _ PostDoc) error { return nil }
func (s *pgFtsSearcher) DeletePost(_ context.Context, _ int64) error  { return nil }
