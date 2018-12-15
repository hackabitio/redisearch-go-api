package handler

import (
	"net/http"
	"io"
	"encoding/json"

	"github.com/7kmCo/redisearch-go-api/redisearch"
)


// Drop a certain index
func (h handler) DelDoc(w io.Writer, r *http.Request) (interface{}, int, error) {
	// First, we need to decode post body from the request
	data := &redisearch.DelDoc{}
	json.NewDecoder(r.Body).Decode(&data)
	// Set index name to the client
	h.client.IndexName(data.IndexName)
	// And finally delete the document
	err := h.client.Delete( data )
	
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	
	return "Document deleted!", http.StatusOK, nil
}