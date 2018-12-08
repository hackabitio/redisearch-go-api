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
	Flags []string `json:"flags"`
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
	// Create search query and set flags
	
	var queryFlags redisearch.Flag = 0x0

	if data.Flags != nil {
		flags := map[string]redisearch.Flag{
			"VERBATIM": 0x1,
			"NOCONTENT": 0x2,
			"WITHSCORES": 0x4,
			"INORDER": 0x08,
			"WITHPAYLOADS": 0x10,
			"NOSTOPWORDS": 0x20,
		}
		for _, f := range data.Flags {
			queryFlags += flags[f]
		}
	}

	query := redisearch.NewQuery(data.Query).Limit(data.From, data.Offset).SetFlags(queryFlags)
	// Then we do the serach
	docs, total, err := client.Search(query)
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