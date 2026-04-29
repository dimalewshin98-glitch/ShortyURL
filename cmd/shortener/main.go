package main

import (
	"net/http"

	"github.com/dimalewshin98-glitch/ShortyURL/internal/config"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/handler"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/logger"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/repository"
	"go.uber.org/zap"
)

func main() {
	cfg := config.NewConfig()
	var repo repository.RepositoryInterface
	var err error
	if err := logger.Initialize(cfg.LogLevel); err != nil {
		panic(err)
	}
	switch cfg.RepositoryType {
	case "file":
		repo, err = repository.NewfileRepository(cfg.FileStoragePath)
	case "memory":
		repo = repository.NewInmemoryRepository()
	case "db":
		repo, err = repository.NewDBRepository(cfg.DatabaseDsn)
	}
	logger.Log.Info("Repository type set to", zap.String("type", cfg.RepositoryType))
	if err != nil {
		panic(err)
	}
	app := NewApp(repo, *cfg)
	appHandler := app.GetHandler()
	logger.Log.Info("Running server", zap.String("address", cfg.ServerHostPort))
	err = http.ListenAndServe(cfg.ServerHostPort, logger.RequestLogger(handler.AuthMiddleware(handler.GzipMiddleware(appHandler), repo)))
	if err != nil {
		logger.Log.Fatal("Server failed", zap.Error(err))
		panic(err)
	}
}
