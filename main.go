package main

import (
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/7kmCo/redisearch-go-api/redisearch"
	"github.com/go-chi/chi"
)

var client *redisearch.Client

func main() {
	client = redisearch.NewClient("localhost:6379", "")
	router := chi.NewRouter()
	router.Get("/info/{idx}", infoHandler)
	router.Post("/search", searchHandler)
	http.ListenAndServe(":8080", router)
}

type SearchQuery struct {
	IndexName string `json:"indexName"`
	Query string `json:"query"`
	From int `json:"from"`
	Offset int `json:"offset"`
}

type SearchResponse struct {
	Total int `json:"total"`
	Results []redisearch.Document `json:"results"`
}

// Handler for search on the index
func searchHandler(w http.ResponseWriter, r *http.Request){
	// First, we need to decode post body from the request
	data := &SearchQuery{}
	json.NewDecoder(r.Body).Decode(&data)
	// We set index name
  client.IndexName(data.IndexName)
	// Then we do the serach
	docs, total, err := client.Search(redisearch.NewQuery(data.Query).Limit(data.From, data.Offset))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// If no results, just return an empty map
	if total == 0 {
		w.Write([]byte{})
		return
	}

	searchResponse := &SearchResponse{
		Total: total,
		Results: docs,
	}
	response, err := json.Marshal(searchResponse)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
	}
	
	fmt.Println(docs[0].Id, docs[0].Properties["post_title"], total, err)

	w.Write(response)
}

// Get index info
func infoHandler(w http.ResponseWriter, r *http.Request){
	// Get index name from url params
	indexName := chi.URLParam(r, "idx")
	// Set index name to the client
	client.IndexName(indexName)
	// Get index info. This converts to "FT.INFO <index_name>" command
	info, err := client.Info()
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	indexInfo, err := json.Marshal(info)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
	}
	
	w.Write(indexInfo)
}