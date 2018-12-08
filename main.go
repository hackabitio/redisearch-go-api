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
	router.Post("/search", searchHandler)
	http.ListenAndServe(":8080", router)
}

type SearchQuery struct {
	IndexName string `json:"indexName"`
	Query string `json:"query"`
	From int `json:"from"`
	Offset int `json:"offset"`
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

	response, err := json.Marshal(docs)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
	}
	
	fmt.Println(docs[0].Id, docs[0].Properties["post_title"], total, err)

	w.Write(response)
}
