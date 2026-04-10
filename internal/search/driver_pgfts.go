package search

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
)

// pgFtsSearcher uses PostgreSQL full-text search (tsvector / to_tsquery).
// Requires PostgreSQL — will NOT work with SQLite.
// Supports Chinese via pg_jieba or zhparser extensions (optional).
type pgFtsSearcher struct{}

func newPgFtsSearcher() Searcher { return &pgFtsSearcher{} }

func (s *pgFtsSearcher) Search(ctx context.Context, ct ContentType, keyword string) (SearchResult, error) {
	meta := metaMap[ct]
	type Row struct {
		Id int64 `orm:"id"`
	}
	var rows []Row
	// plainto_tsquery handles multi-word input safely (no need to escape).
	// Uses the "simple" dictionary so it works without language-specific configs.
	err := g.DB().Ctx(ctx).
		Model(meta.Table).
		Fields("id").
		WhereNull("deleted_at").
		Where(fmt.Sprintf("to_tsvector('simple', %s) @@ plainto_tsquery('simple', ?)", meta.FtsCols), keyword).
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

func (s *pgFtsSearcher) Index(_ context.Context, _ ContentType, _ Document) error { return nil }
func (s *pgFtsSearcher) Delete(_ context.Context, _ ContentType, _ int64) error  { return nil }
