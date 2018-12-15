package handler

import (
	"net/http"
	"io"
	"fmt"
	"encoding/json"

	"github.com/7kmCo/redisearch-go-api/redisearch"
	"github.com/go-chi/chi"	
	"github.com/go-chi/cors"
)

type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"response"`
}

type handler struct {
	client *redisearch.Client
}

func New(client *redisearch.Client) chi.Router {
	h := handler{client}
	router := chi.NewRouter()

	// Basic CORS
  cors := cors.New(cors.Options{
    // AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
    AllowedOrigins:   []string{"*"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
    ExposedHeaders:   []string{"Link"},
    AllowCredentials: true,
    MaxAge:           300, // Maximum value not ignored by any of major browsers
		Debug:            true,
  })
	router.Use(cors.Handler)
	
	router.Get("/info/{idx}", requestHandler(h.Info))
	router.Post("/explain", requestHandler(h.Explain))
	router.Post("/search", requestHandler(h.Search))
	router.Get("/get/{idx}/{docId}", requestHandler(h.Get))
	router.Post("/create", requestHandler(h.Create))
	router.Post("/add", requestHandler(h.Add))
	router.Delete("/drop/{idx}", requestHandler(h.Drop))
	router.Delete("/delete", requestHandler(h.DelDoc))
	router.Route("/suggestion", func(r chi.Router) {
		r.Get("/length/{idx}", requestHandler(h.SugLen))
		r.Post("/add", requestHandler(h.SugAdd))
		r.Post("/get", requestHandler(h.SugGet))
		r.Delete("/delete", requestHandler(h.SugDel))
	})
	// Router group for synonyms
	router.Route("/synonym", func(r chi.Router) {
		r.Post("/add", requestHandler(h.SynAdd))
	})
	http.ListenAndServe(":8080", router)
	return router
}

func requestHandler(h func(io.Writer, *http.Request) (interface{}, int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, status, err := h(w, r)
		if err != nil {
			data = err.Error()
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		fmt.Println(data)

		err = json.NewEncoder(w).Encode(response{Data: data, Success: err == nil})
		if err != nil {
			fmt.Printf("could not encode response to output: %v", err)
		}
	}
}

