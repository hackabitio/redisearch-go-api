package main

import (
	"fmt"
	"github.com/7kmCo/redisearch-go-api/redisearch"
)

func main() {
	c := redisearch.NewClient("localhost:6379", "myIndex")
	docs, total, err := c.Search(redisearch.NewQuery("lorem").Limit(0, 2))

	fmt.Println(docs[0].Id, docs[0].Properties["post_title"], total, err)
}