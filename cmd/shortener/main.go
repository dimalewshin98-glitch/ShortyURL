package main

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/dimalewshin98-glitch/ShortyURL/internal/config"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/handler"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/logger"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/repository"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.NewConfig()
	if err := logger.Initialize(cfg.LogLevel); err != nil {
		panic(err)
	}
	repository := repository.NewInmemoryRepository()
	shorterService := service.NewShorterService(repository, cfg)
	requestsHandler := handler.NewRequestsHandler(shorterService)
	r := chi.NewRouter()
	r.Post("/", logger.RequestLogger(requestsHandler.Shorten))
	r.Post("/api/shorten", logger.RequestLogger(requestsHandler.ApiShorten))
	r.Get("/{id}", logger.RequestLogger(requestsHandler.GetURL))
	logger.Log.Info("Running server", zap.String("address", cfg.ServerHostPort))
	err := http.ListenAndServe(cfg.ServerHostPort, r)
	if err != nil {
		logger.Log.Fatal("Server failed", zap.Error(err))
	}
}
