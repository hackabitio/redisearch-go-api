package handler

import (
	"net/http"
	"io"

	"github.com/go-chi/chi"	
)

// Drop a certain index
func (h handler) Drop(w io.Writer, r *http.Request) (interface{}, int, error) {
	// Get index name from url params
	indexName := chi.URLParam(r, "idx")
	// Set index name to the client
	h.client.IndexName(indexName)
	// Simply drop the index
	err := h.client.Drop()
	
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	
	return "Index dropped!", http.StatusOK, nil
}