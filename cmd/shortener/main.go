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
	var repo repository.RepositoryInterface
	var err error
	switch cfg.RepositoryType {
	case "file":
		repo, err = repository.NewfileRepository(cfg.FileStoragePath)
	case "memory":
		repo = repository.NewInmemoryRepository()
	}
	if err != nil {
		panic(err)
	}
	if repo == nil {
		panic("Repo init error")
	}
	shorterService := service.NewShorterService(repo, cfg)
	requestsHandler := handler.NewRequestsHandler(shorterService)
	r := chi.NewRouter()
	r.Post("/", requestsHandler.Shorten)
	r.Get("/{id}", requestsHandler.GetURL)
	r.Post("/api/shorten", requestsHandler.ApiShorten)
	logger.Log.Info("Running server", zap.String("address", cfg.ServerHostPort))
	err = http.ListenAndServe(cfg.ServerHostPort, logger.RequestLogger(handler.GzipMiddleware(r)))
	if err != nil {
		logger.Log.Fatal("Server failed", zap.Error(err))
	}
}
