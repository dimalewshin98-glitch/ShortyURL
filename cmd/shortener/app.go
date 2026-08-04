package main

import (
	"net/http"

	"github.com/dimalewshin98-glitch/ShortyURL/internal/config"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/handler"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/repository"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/service"
	"github.com/go-chi/chi/v5"
)

type App struct {
	repo           repository.RepositoryInterface
	cfg            config.Config
	shorterService *service.ShorterService
}

func NewApp(repo repository.RepositoryInterface, cfg config.Config, auditors []service.Auditor) *App {
	shorterService := service.NewShorterService(repo, &cfg)
	for _, a := range auditors {
		shorterService.AttachAuditor(a)
	}
	return &App{
		repo:           repo,
		cfg:            cfg,
		shorterService: shorterService,
	}
}

func (a *App) GetHTTPHandler() http.Handler {
	requestsHandler := handler.NewRequestsHandler(a.shorterService)
	r := chi.NewRouter()
	r.Get("/ping", requestsHandler.Ping)
	r.Get("/api/internal/stats", requestsHandler.ApiInternalStats)
	r.Post("/", requestsHandler.Shorten)
	r.Get("/{id}", requestsHandler.GetURL)
	r.Post("/api/shorten", requestsHandler.ApiShorten)
	r.Post("/api/shorten/batch", requestsHandler.ApiShortenBatch)
	r.Get("/api/user/urls", requestsHandler.ApiUserUrls)
	r.Delete("/api/user/urls", requestsHandler.Delete)
	return r
}

func (a *App) GetGRPCRequestsHandler() *handler.GRPCRequestsHandler {
	GRPCRequestsHandler := handler.NewGRPCRequestsHandler(a.shorterService)
	return GRPCRequestsHandler
}
