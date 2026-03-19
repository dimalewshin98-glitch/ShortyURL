package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dimalewshin98-glitch/ShortyURL/internal/handler"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/repository"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	repository := repository.NewInmemoryRepository()
	shorterService := service.NewShorterService(repository)
	requestsHandler := handler.NewRequestsHandler(shorterService)
	r := chi.NewRouter()
	r.Post("/", requestsHandler.Shorten)
	r.Get("/{id}", requestsHandler.GetUrl)
	fmt.Println("server starting at port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
