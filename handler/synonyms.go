package handler

import (
	"net/http"
	"io"
	"encoding/json"
)


// Add synonym group
func (h handler) SynAdd(w io.Writer, r *http.Request) (interface{}, int, error) {
	// First, decode post body from the request
	data := &struct{
		IndexName 	string	`json:"indexName"`
		Term 				string	`json:"term"`
		Synonym			string `json:"synonym"`
	}{}
	json.NewDecoder(r.Body).Decode(&data)
	// If suggestion index name is not available in the request, return error
	if data.IndexName == "" {
		return nil, http.StatusInternalServerError, nil
	}
	// Set index name to the client
	h.client.IndexName(data.IndexName)
	
	_, err := h.client.SynAdd(data.Term, data.Synonym)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	
	return "Synonym added", http.StatusOK, nil
}
