# RediSearch Go Client

Go client for [RediSearch](http://redisearch.io), based on redigo.

# Installing 

```sh
go get https://github.com/7kmCo/redisearch-go-api/redisearch
```

# Usage Example

```go

import (
	"fmt"
	"log"
	"time"

	"https://github.com/7kmCo/redisearch-go-api/redisearch"
)

func ExampleClient() {

	// Create a client. By default a client is schemaless
	// unless a schema is provided when creating the index
	c := redisearch.NewClient("localhost:6379", "myIndex")

	// Create a schema
	sc := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewTextField("body")).
		AddField(redisearch.NewTextFieldOptions("title", redisearch.TextFieldOptions{Weight: 5.0, Sortable: true})).
		AddField(redisearch.NewNumericField("date"))

	// Drop an existing index. If the index does not exist an error is returned
	c.Drop()

	// Create the index with the given schema
	if err := c.CreateIndex(sc); err != nil {
		log.Fatal(err)
	}

	// Create a document with an id and given score
	doc := redisearch.NewDocument("doc1", 1.0)
	doc.Set("title", "Hello world").
		Set("body", "foo bar").
		Set("date", time.Now().Unix())

	// Index the document. The API accepts multiple documents at a time
	if err := c.Index([]redisearch.Document{doc}...); err != nil {
		log.Fatal(err)
	}

	// Searching with limit and sorting
	docs, total, err := c.Search(redisearch.NewQuery("hello world").
		Limit(0, 2).
		SetReturnFields("title"))

	fmt.Println(docs[0].Id, docs[0].Properties["title"], total, err)
	// Output: doc1 Hello world 1 <nil>
}
```

**Note**: This package is cloned from [redisearch-go](https://github.com/RedisLabs/redisearch-go) created by RedisLabs. The reason I didn't just used thier package was that the repo was not active recently, its dependencies was outdated and I needed to do some modifications. So just cloned it and did my edits.