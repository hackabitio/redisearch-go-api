package handler

import (
	"net/http"
	"io"
	// "fmt"
	"encoding/json"

	"github.com/7kmCo/redisearch-go-api/redisearch"
)

// SpellCheck, Performs spelling correction on a query,
// returning suggestions for misspelled terms
func (h handler) Check(w io.Writer, r *http.Request) (interface{}, int, error) {
	// First, decode post body from the request
	data := &SearchQuery{}
	json.NewDecoder(r.Body).Decode(&data)
	// If suggestion index name is not available in the request, return error
	if data.IndexName == "" {
		return nil, http.StatusInternalServerError, nil
	}
	// Set index name to the client
	h.client.IndexName(data.IndexName)
	// Create new query
	query := redisearch.NewQuery(data.Query)
	// Then we do the SpellCheck
	spelled, err := h.client.SpellCheck(query)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	
	return spelled, http.StatusOK, nil
}
