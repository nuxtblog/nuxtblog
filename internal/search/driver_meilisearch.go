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
//	    host:   "http://localhost:7700"
//	    apiKey: ""          # master key or search-only key
//	    index:  "posts"     # index name (created automatically on first IndexPost)
type meiliSearcher struct {
	host   string
	apiKey string
	index  string
	client *http.Client
}

func newMeiliSearcher(ctx context.Context) Searcher {
	host := g.Cfg().MustGet(ctx, "search.meilisearch.host", "http://localhost:7700").String()
	apiKey := g.Cfg().MustGet(ctx, "search.meilisearch.apiKey", "").String()
	index := g.Cfg().MustGet(ctx, "search.meilisearch.index", "posts").String()
	return &meiliSearcher{
		host:   host,
		apiKey: apiKey,
		index:  index,
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

func (s *meiliSearcher) SearchPostIDs(_ context.Context, keyword string) ([]int64, error) {
	type searchReq struct {
		Q                  string   `json:"q"`
		AttributesToRetrieve []string `json:"attributesToRetrieve"`
		Limit              int      `json:"limit"`
	}
	type hit struct {
		ID int64 `json:"id"`
	}
	type searchRes struct {
		Hits []hit `json:"hits"`
	}

	resp, err := s.do("POST", fmt.Sprintf("/indexes/%s/search", s.index), searchReq{
		Q:                  keyword,
		AttributesToRetrieve: []string{"id"},
		Limit:              1000,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result searchRes
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	ids := make([]int64, len(result.Hits))
	for i, h := range result.Hits {
		ids[i] = h.ID
	}
	return ids, nil
}

func (s *meiliSearcher) IndexPost(_ context.Context, doc PostDoc) error {
	_, err := s.do("POST", fmt.Sprintf("/indexes/%s/documents", s.index), []PostDoc{doc})
	return err
}

func (s *meiliSearcher) DeletePost(_ context.Context, id int64) error {
	_, err := s.do("DELETE", fmt.Sprintf("/indexes/%s/documents/%d", s.index, id), nil)
	return err
}
