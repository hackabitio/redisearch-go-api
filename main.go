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
	router.Post("/create", createHandler)
	router.Post("/add", addHandler)
	router.Delete("/drop/{idx}", dropHandler)
	router.Route("/suggestion", func(r chi.Router) {
		r.Post("/add", addSuggestion)
	})
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

// Field represents a single field's Schema
type Field struct {
	Name     string `json:"name"`
	Type     string	`json:"type"`
	Sortable bool 	`json:"sortable"`
	Options  map[string]interface{} `json:"options"`
}

type NewIndex struct {
	IndexName string `json:"indexName"`
	Schema []Field `json:"schema"`
}

// Handler for creating index
func createHandler(w http.ResponseWriter, r *http.Request){
	// First, decode post body from the request
	data := &NewIndex{}
	json.NewDecoder(r.Body).Decode(&data)
	// We set index name
	client.IndexName(data.IndexName)
	// Create schema
	if data.Schema == nil {
		http.Error(w, "Request doesn't include schema", http.StatusBadRequest)
		return
	}
	// Initiate schema
	newSchema := redisearch.NewSchema(redisearch.DefaultOptions)

	// Then we iterate over schema received from request
	for _, sc := range data.Schema {
		// We must handle different types separately, so we can sanitize fields and do the validation stuff
		switch sc.Type {

			// Text field
			case "text":
				// If any options provided, do the validations here
				if sc.Options != nil {
					var textFieldOptions redisearch.TextFieldOptions
					if sc.Options["weight"] != nil{
						textFieldOptions.Weight = float32(sc.Options["weight"].(float64))
					}
					if sc.Options["sortable"] != nil {
						textFieldOptions.Sortable = sc.Options["sortable"].(bool)
					}
					if sc.Options["noStem"] != nil {
						textFieldOptions.NoStem = sc.Options["noStem"].(bool)
					}
					if sc.Options["noIndex"] != nil {
						textFieldOptions.NoIndex = sc.Options["noIndex"].(bool)
					}
					newSchema.AddField(redisearch.NewTextFieldOptions(sc.Name, textFieldOptions))
				} else {
					// If there are no options with field schema,
					// simply use its name to create the schema with default options
					newSchema.AddField(redisearch.NewTextField(sc.Name))
				}

			// Numeric field
			case "numeric":
				// If any options provided, do the validations here
				if sc.Options != nil {
					var numFieldOptions redisearch.NumericFieldOptions
					if sc.Options["sortable"] != nil {
						numFieldOptions.Sortable = sc.Options["sortable"].(bool)
					}
					if sc.Options["noIndex"] != nil {
						numFieldOptions.NoIndex = sc.Options["noIndex"].(bool)
					}
					newSchema.AddField(redisearch.NewNumericFieldOptions(sc.Name, numFieldOptions))
				} else {
					// numeric field with default options
					newSchema.AddField(redisearch.NewNumericField(sc.Name))
				}


			// Tag field
			case "tag":
				// If any options provided, do the validations here
				if sc.Options != nil {
					var tagFieldOptions redisearch.TagFieldOptions
					if sc.Options["separator"] != nil{
						tagFieldOptions.Separator = []byte(sc.Options["separator"].(string))[0]
					}
					if sc.Options["sortable"] != nil {
						tagFieldOptions.Sortable = sc.Options["sortable"].(bool)
					}
					if sc.Options["noIndex"] != nil {
						tagFieldOptions.NoIndex = sc.Options["noIndex"].(bool)
					}
					newSchema.AddField(redisearch.NewTagFieldOptions(sc.Name, tagFieldOptions))
				} else {
					// Tag field with default options
					newSchema.AddField(redisearch.NewTagField(sc.Name))
				}
		}

	}
	if err := client.CreateIndex(newSchema); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	fmt.Println(newSchema)

	// w.Write(response)
}

type Document struct {
	IndexName string `json:"indexName"`
	Doc redisearch.Document `json:"document"`
}

// Add document
func addHandler(w http.ResponseWriter, r *http.Request){
	// First, decode post body from the request
	data := &Document{}
	json.NewDecoder(r.Body).Decode(&data)
	// Set index name to the client
	client.IndexName(data.IndexName)

	// Create a document with an id, given score and its fields
	document := redisearch.NewDocument(data.Doc.Id, data.Doc.Score)
	for f, v := range data.Doc.Properties {
		document.Set(f, v)
	}

	// Index the document
	if err := client.Index([]redisearch.Document{document}...); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	docAdded, _ := json.Marshal("Document added to the index")
	
	w.Write(docAdded)
}

// Drop a certain index
func dropHandler(w http.ResponseWriter, r *http.Request){
	// Get index name from url params
	indexName := chi.URLParam(r, "idx")
	// Set index name to the client
	client.IndexName(indexName)
	// Simply drop the index
	err := client.Drop()
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	dropped, _ := json.Marshal("Index successfully dropped")
	
	w.Write(dropped)
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

type Suggestion struct {
	IndexName string `json:"name"`
	Sugg redisearch.Suggestion `json:"suggestion"`
}
// Add suggestion/autocomplete
func addSuggestion(w http.ResponseWriter, r *http.Request){
	// First, decode post body from the request
	data := &Suggestion{}
	json.NewDecoder(r.Body).Decode(&data)
	// Create new connection pool
	autocompleter := redisearch.NewAutocompleter("localhost:6379", "")
	// If suggestion index name is not available in the request, return error
	if data.IndexName == "" {
		http.Error(w, "Suggestion index name is not provided", http.StatusBadRequest)
		return
	}
	// Set index name to the client
	autocompleter.IndexName(data.IndexName)
	
	err := autocompleter.AddTerms(data.Sugg)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	added, _ := json.Marshal("Suggestion added")
	
	w.Write(added)
}