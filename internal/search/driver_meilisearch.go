package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// meiliSearcher integrates with MeiliSearch.
// Config (config.yaml):
//
//	search:
//	  driver: "meilisearch"
//	  meilisearch:
//	    host:        "http://localhost:7700"
//	    apiKey:      ""
//	    index:       "posts"
//	    docIndex:    "docs"
//	    momentIndex: "moments"
type meiliSearcher struct {
	host    string
	apiKey  string
	indexes map[ContentType]string
	client  *http.Client
}

func newMeiliSearcher(ctx context.Context) Searcher {
	host := g.Cfg().MustGet(ctx, "search.meilisearch.host", "http://localhost:7700").String()
	apiKey := g.Cfg().MustGet(ctx, "search.meilisearch.apiKey", "").String()
	return &meiliSearcher{
		host:   host,
		apiKey: apiKey,
		indexes: map[ContentType]string{
			ContentPost:   g.Cfg().MustGet(ctx, "search.meilisearch.index", "posts").String(),
			ContentDoc:    g.Cfg().MustGet(ctx, "search.meilisearch.docIndex", "docs").String(),
			ContentMoment: g.Cfg().MustGet(ctx, "search.meilisearch.momentIndex", "moments").String(),
		},
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (s *meiliSearcher) do(method, path string, body any) (*http.Response, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, s.host+path, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}
	return s.client.Do(req)
}

func (s *meiliSearcher) Search(_ context.Context, ct ContentType, keyword string) (SearchResult, error) {
	idx := s.indexes[ct]
	type searchReq struct {
		Q                    string   `json:"q"`
		AttributesToRetrieve []string `json:"attributesToRetrieve"`
		Limit                int      `json:"limit"`
	}
	type hit struct {
		ID int64 `json:"id"`
	}
	type searchRes struct {
		Hits               []hit `json:"hits"`
		EstimatedTotalHits int   `json:"estimatedTotalHits"`
	}

	resp, err := s.do("POST", fmt.Sprintf("/indexes/%s/search", idx), searchReq{
		Q:                    keyword,
		AttributesToRetrieve: []string{"id"},
		Limit:                defaultSearchLimit,
	})
	if err != nil {
		return SearchResult{}, err
	}
	defer resp.Body.Close()

	var result searchRes
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return SearchResult{}, err
	}
	ids := make([]int64, len(result.Hits))
	for i, h := range result.Hits {
		ids[i] = h.ID
	}
	return SearchResult{IDs: ids, Total: result.EstimatedTotalHits}, nil
}

func (s *meiliSearcher) Index(_ context.Context, ct ContentType, doc Document) error {
	idx := s.indexes[ct]
	_, err := s.do("POST", fmt.Sprintf("/indexes/%s/documents", idx), []Document{doc})
	return err
}

func (s *meiliSearcher) Delete(_ context.Context, ct ContentType, id int64) error {
	idx := s.indexes[ct]
	_, err := s.do("DELETE", fmt.Sprintf("/indexes/%s/documents/%d", idx, id), nil)
	return err
}
