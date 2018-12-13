package handler

import (
	"net/http"
	"io"
	"encoding/json"

	"github.com/7kmCo/redisearch-go-api/redisearch"
)

type Suggestion struct {
	IndexName string `json:"name"`
	Sugg []redisearch.Suggestion `json:"suggestion"`
}

// Add suggestion/autocomplete
func (h handler) SugAdd(w io.Writer, r *http.Request) (interface{}, int, error) {
	// First, decode post body from the request
	data := &Suggestion{}
	json.NewDecoder(r.Body).Decode(&data)
	// Create new connection pool
	autocompleter := redisearch.NewAutocompleter("localhost:6379", "")
	// If suggestion index name is not available in the request, return error
	if data.IndexName == "" {
		return nil, http.StatusInternalServerError, nil
	}
	// Set index name to the client
	autocompleter.IndexName(data.IndexName)
	
	err := autocompleter.AddTerms(data.Sugg ...)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	
	return "Suggestion added", http.StatusOK, nil
}

type Suggest struct {
	IndexName string `json:"name"`
	Prefix string `json:"prefix"`
	Options redisearch.SuggestOptions `json:"options"`
}

// Add suggestion/autocomplete
func (h handler) SugGet(w io.Writer, r *http.Request) (interface{}, int, error) {
	// First, decode post body from the request
	data := &Suggest{}
	json.NewDecoder(r.Body).Decode(&data)
	// Create new connection pool
	autocompleter := redisearch.NewAutocompleter("localhost:6379", "")
	// If suggestion index name is not available in the request, return error
	if data.IndexName == "" {
		return nil, http.StatusInternalServerError, nil
	}
	// Set index name to the client
	autocompleter.IndexName(data.IndexName)
	// Then get suggestion
	sugg, err := autocompleter.SuggestOpts(data.Prefix, data.Options)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	
	return sugg, http.StatusOK, nil
}

// Deletes a string from suggestion/autocomplete
func (h handler) SugDel(w io.Writer, r *http.Request) (interface{}, int, error) {
	// First, decode post body from the request
	data := &struct{
		IndexName 	string	`json:"name"`
		Prefix 			string	`json:"text"`
	}{}
	json.NewDecoder(r.Body).Decode(&data)
	// Create new connection pool
	autocompleter := redisearch.NewAutocompleter("localhost:6379", "")
	// If suggestion index name is not available in the request, return error
	if data.IndexName == "" {
		return nil, http.StatusInternalServerError, nil
	}
	// Set index name to the client
	autocompleter.IndexName(data.IndexName)
	// Delete suggestion from index
	err := autocompleter.Delete(data.Prefix)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	
	return "String deleted", http.StatusOK, nil
}
