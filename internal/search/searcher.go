// Package search provides a pluggable search abstraction.
// Switch drivers by setting search.driver in config.yaml:
//
//	search:
//	  driver: "db"           # default — SQL LIKE, works with SQLite & PostgreSQL
//	  driver: "postgres_fts" # PostgreSQL full-text search (tsvector)
//	  driver: "meilisearch"  # MeiliSearch
package search

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

// ContentType identifies the type of searchable content.
type ContentType string

const (
	ContentPost   ContentType = "post"
	ContentDoc    ContentType = "doc"
	ContentMoment ContentType = "moment"
)

// Document is the unified struct for external search indexes.
type Document struct {
	ID      int64
	Title   string
	Excerpt string
	Content string
}

// PostDoc and DocDoc are aliases kept for backward compatibility.
type PostDoc = Document
type DocDoc = Document

// SearchResult holds search IDs and the total match count.
// If Total > len(IDs), results were truncated by the search limit.
type SearchResult struct {
	IDs   []int64
	Total int // total matches (before limit); -1 if unknown
}

// contentMeta holds per-type DB metadata.
type contentMeta struct {
	Table   string // DB table name
	LikeCol string // column for SQL LIKE
	FtsCols string // tsvector expression
}

var metaMap = map[ContentType]contentMeta{
	ContentPost:   {Table: "posts", LikeCol: "title", FtsCols: "title || ' ' || excerpt"},
	ContentDoc:    {Table: "docs", LikeCol: "title", FtsCols: "title || ' ' || excerpt"},
	ContentMoment: {Table: "moments", LikeCol: "content", FtsCols: "content"},
}

const defaultSearchLimit = 1000

// Searcher is the pluggable search interface.
// All implementations must be safe for concurrent use.
type Searcher interface {
	// Search returns IDs matching keyword for the given content type.
	Search(ctx context.Context, ct ContentType, keyword string) (SearchResult, error)

	// Index upserts a document into the search index.
	// No-op for the "db" and "postgres_fts" drivers.
	Index(ctx context.Context, ct ContentType, doc Document) error

	// Delete removes a document from the search index.
	// No-op for the "db" and "postgres_fts" drivers.
	Delete(ctx context.Context, ct ContentType, id int64) error
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
