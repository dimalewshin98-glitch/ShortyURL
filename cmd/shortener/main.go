package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dimalewshin98-glitch/ShortyURL/internal/config"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/handler"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/repository"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	config := config.NewConfig()
	repository := repository.NewInmemoryRepository()
	shorterService := service.NewShorterService(repository, config)
	requestsHandler := handler.NewRequestsHandler(shorterService)
	r := chi.NewRouter()
	r.Post("/", requestsHandler.Shorten)
	r.Get("/{id}", requestsHandler.GetURL)
	fmt.Println("server starting at: " + config.ServerHostPort)
	log.Fatal(http.ListenAndServe(config.ServerHostPort, r))
}
