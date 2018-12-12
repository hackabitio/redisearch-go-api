package handler

import (
	"net/http"
	"io"
	"encoding/json"

	"github.com/7kmCo/redisearch-go-api/redisearch"
)

type Document struct {
	IndexName string `json:"indexName"`
	Doc redisearch.Document `json:"document"`
}

// Add document
func (h handler) Add(w io.Writer, r *http.Request) (interface{}, int, error) {
	// First, decode post body from the request
	data := &Document{}
	json.NewDecoder(r.Body).Decode(&data)
	// Set index name to the client
	h.client.IndexName(data.IndexName)

	// Create a document with an id, given score and its fields
	document := redisearch.NewDocument(data.Doc.Id, data.Doc.Score)
	for f, v := range data.Doc.Properties {
		document.Set(f, v)
	}

	// Index the document
	if err := h.client.Index([]redisearch.Document{document}...); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	
	return "Document added to the index", http.StatusOK, nil
}
