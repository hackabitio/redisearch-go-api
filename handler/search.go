package handler

import (
	"net/http"
	"io"
	"fmt"
	"encoding/json"

	"github.com/7kmCo/redisearch-go-api/redisearch"
)

type SearchQuery struct {
	IndexName string `json:"indexName"`
	Query string `json:"query"`
	Flags []string `json:"flags"`
	From int `json:"from"`
	Offset int `json:"offset"`
}

type SearchResponse struct {
	Total int `json:"total"`
	Results []redisearch.Document `json:"results"`
}

// Handler for search on the index
func (h handler) Search(w io.Writer, r *http.Request) (interface{}, int, error) {
	// First, we need to decode post body from the request
	data := &SearchQuery{}
	json.NewDecoder(r.Body).Decode(&data)
	// We set index name
	h.client.IndexName(data.IndexName)
	// Create search query and set flags
	
	var queryFlags redisearch.Flag = 0x0

	if data.Flags != nil {
		flags := map[string]redisearch.Flag{
			"VERBATIM": 0x1,
			"NOCONTENT": 0x2,
			"WITHSCORES": 0x4,
			"INORDER": 0x08,
			"WITHPAYLOADS": 0x10,
			"NOSTOPWORDS": 0x20,
		}
		for _, f := range data.Flags {
			queryFlags += flags[f]
		}
	}

	query := redisearch.NewQuery(data.Query).Limit(data.From, data.Offset).SetFlags(queryFlags)
	// Then we do the serach
	docs, total, err := h.client.Search(query)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	// If no results, just return an empty map
	if total == 0 {
		return nil, http.StatusInternalServerError, err
	}

	searchResponse := &SearchResponse{
		Total: total,
		Results: docs,
	}
	
	fmt.Println(docs[0].Id, docs[0].Properties["post_title"], total, err)

	return searchResponse, http.StatusOK, nil
}