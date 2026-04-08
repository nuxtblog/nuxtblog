// Package search provides a pluggable search abstraction.
// Switch drivers by setting search.driver in config.yaml:
//   search:
//     driver: "db"           # default — SQL LIKE, works with SQLite & PostgreSQL
//     driver: "postgres_fts" # PostgreSQL full-text search (tsvector)
//     driver: "meilisearch"  # MeiliSearch
package search

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

// PostDoc is the document pushed to the search index.
type PostDoc struct {
	ID      int64
	Title   string
	Excerpt string
	Content string
}

// Searcher is the pluggable search interface.
// All implementations must be safe for concurrent use.
type Searcher interface {
	// SearchPostIDs returns IDs of posts matching keyword.
	// Returns all matches (caller handles pagination via SQL WHERE id IN).
	SearchPostIDs(ctx context.Context, keyword string) ([]int64, error)

	// IndexPost upserts a post document into the search index.
	// No-op for the "db" driver (DB is the index).
	IndexPost(ctx context.Context, doc PostDoc) error

	// DeletePost removes a post from the search index.
	// No-op for the "db" driver.
	DeletePost(ctx context.Context, id int64) error
}

var defaultSearcher Searcher

// Default returns the singleton Searcher configured in config.yaml.
// Falls back to the "db" driver if not configured.
func Default(ctx context.Context) Searcher {
	if defaultSearcher != nil {
		return defaultSearcher
	}
	driver := g.Cfg().MustGet(ctx, "search.driver", "db").String()
	switch driver {
	case "meilisearch":
		defaultSearcher = newMeiliSearcher(ctx)
	case "postgres_fts":
		defaultSearcher = newPgFtsSearcher()
	default:
		defaultSearcher = newDbSearcher()
	}
	return defaultSearcher
}

// Reset clears the cached singleton (useful for tests).
func Reset() { defaultSearcher = nil }
