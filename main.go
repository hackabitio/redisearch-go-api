package main

import (
	"fmt"
	"net/http"

	"github.com/7kmCo/redisearch-go-api/redisearch"
	"github.com/7kmCo/redisearch-go-api/handler"
)

func main() {
	client := redisearch.NewClient("localhost:6379", "")

	// Create a server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", "localhost", "8080"),
		Handler: handler.New(client),
	}

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("%v", err)
	} else {
		fmt.Println("Server closed!")
	}

}

// type Suggestion struct {
// 	IndexName string `json:"name"`
// 	Sugg redisearch.Suggestion `json:"suggestion"`
// }
// // Add suggestion/autocomplete
// func addSuggestion(w http.ResponseWriter, r *http.Request){
// 	// First, decode post body from the request
// 	data := &Suggestion{}
// 	json.NewDecoder(r.Body).Decode(&data)
// 	// Create new connection pool
// 	autocompleter := redisearch.NewAutocompleter("localhost:6379", "")
// 	// If suggestion index name is not available in the request, return error
// 	if data.IndexName == "" {
// 		http.Error(w, "Suggestion index name is not provided", http.StatusBadRequest)
// 		return
// 	}
// 	// Set index name to the client
// 	autocompleter.IndexName(data.IndexName)
	
// 	err := autocompleter.AddTerms(data.Sugg)

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
	
// 	added, _ := json.Marshal("Suggestion added")
	
// 	w.Write(added)
// }

// type Suggest struct {
// 	IndexName string `json:"name"`
// 	Prefix string `json:"prefix"`
// 	Options redisearch.SuggestOptions `json:"options"`
// }
// // Add suggestion/autocomplete
// func getSuggestion(w http.ResponseWriter, r *http.Request){
// 	// First, decode post body from the request
// 	data := &Suggest{}
// 	json.NewDecoder(r.Body).Decode(&data)
// 	// Create new connection pool
// 	autocompleter := redisearch.NewAutocompleter("localhost:6379", "")
// 	// If suggestion index name is not available in the request, return error
// 	if data.IndexName == "" {
// 		http.Error(w, "Suggestion index name is not provided", http.StatusBadRequest)
// 		return
// 	}
// 	// Set index name to the client
// 	autocompleter.IndexName(data.IndexName)
// 	// The get suggestion
// 	sugg, err := autocompleter.SuggestOpts(data.Prefix, data.Options)

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
	
// 	added, _ := json.Marshal(sugg)
	
// 	w.Write(added)
// }