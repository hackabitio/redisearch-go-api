package handler

import (
	"net/http"
	"io"

	"github.com/go-chi/chi"
)

// Get index info
func (h handler) Info(w io.Writer, r *http.Request) (interface{}, int, error) {
	// Get index name from url params
	indexName := chi.URLParam(r, "idx")
	// Set index name to the client
	h.client.IndexName(indexName)
	// Get index info. This converts to "FT.INFO <index_name>" command
	info, err := h.client.Info()
  if err != nil {
    return nil, http.StatusInternalServerError, err
	}
	
	return info, http.StatusOK, nil
}
