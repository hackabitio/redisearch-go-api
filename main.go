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
