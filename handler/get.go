package handler

import (
	"net/http"
	"io"

	"github.com/go-chi/chi"
)

// Returns the full contents of a document.
func (h handler) Get(w io.Writer, r *http.Request) (interface{}, int, error) {
	// Get index name from url params
	indexName := chi.URLParam(r, "idx")
	// Get document id from url params
	docId := chi.URLParam(r, "docId")
	// Set index name to the client
	h.client.IndexName(indexName)
	// Get document
	doc, err := h.client.Get(docId)
  if err != nil {
    return nil, http.StatusInternalServerError, err
	}
	
	return doc, http.StatusOK, nil
}
