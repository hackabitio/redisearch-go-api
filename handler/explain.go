package handler

import (
	"net/http"
	"io"
	"encoding/json"

	"github.com/7kmCo/redisearch-go-api/redisearch"
)

// Explain search query. This is useful for debugging
func (h handler) Explain(w io.Writer, r *http.Request) (interface{}, int, error) {// First, we need to decode post body from the request
	data := &SearchQuery{}
	json.NewDecoder(r.Body).Decode(&data)
	// We set index name
	h.client.IndexName(data.IndexName)
	
	query := redisearch.NewQuery(data.Query)
	// Get index info. This converts to "FT.INFO <index_name>" command
	explain, err := h.client.Explain(query)
  if err != nil {
    return nil, http.StatusInternalServerError, err
	}
	
	return explain, http.StatusOK, nil
}
