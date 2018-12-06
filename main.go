package main

import (
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/7kmCo/redisearch-go-api/redisearch"
	"github.com/go-chi/chi"
)

func main() {
		r := chi.NewRouter()
	r.Get("/search/{query}", searchHandler)
	http.ListenAndServe(":8080", r)
}

func searchHandler(w http.ResponseWriter, r *http.Request){
	query := chi.URLParam(r, "query")
	
	c := redisearch.NewClient("localhost:6379", "myIndex")
	docs, total, err := c.Search(redisearch.NewQuery(query).Limit(0, 10))
	
	if err != nil {
		return
	}
	
	fmt.Println(docs[0].Id, docs[0].Properties["post_title"], total, err)

	js, err := json.Marshal(docs)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
	}
	
	w.Write(js)
	
}