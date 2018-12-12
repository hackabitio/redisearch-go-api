package handler

import (
	"net/http"
	"io"
	"fmt"
	"encoding/json"

	"github.com/7kmCo/redisearch-go-api/redisearch"
)

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
func (h handler) Create(w io.Writer, r *http.Request) (interface{}, int, error) {
	// First, decode post body from the request
	data := &NewIndex{}
	json.NewDecoder(r.Body).Decode(&data)
	// We set index name
	h.client.IndexName(data.IndexName)
	// Create schema
	if data.Schema == nil {
		return nil, http.StatusInternalServerError, nil
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
	if err := h.client.CreateIndex(newSchema); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	fmt.Println(newSchema)

	return newSchema, http.StatusOK, nil
}
