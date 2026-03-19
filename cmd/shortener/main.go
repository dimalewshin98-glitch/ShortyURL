package main

import (
	"fmt"
	"net/http"

	"github.com/dimalewshin98-glitch/ShortyURL/internal/handler"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/repository"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/service"
)

func main() {
	repository := repository.NewInmemoryRepository()
	shorterService := service.NewShorterService(repository)
	requestsHandler := handler.NewRequestsHandler(shorterService)
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, requestsHandler.Shorten)
	mux.HandleFunc(`/{id}`, requestsHandler.GetUrl)
	fmt.Println("server starting at port 8080")
	http.ListenAndServe(`:8080`, mux)
}
