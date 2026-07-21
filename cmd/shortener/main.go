package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/dimalewshin98-glitch/ShortyURL/internal/config"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/handler"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/logger"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/repository"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/service"
	"go.uber.org/zap"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func validateBuildInfo(value string) string {
	if value == "" {
		return "N/A"
	}
	return value
}

func printBuildInfo() {
	fmt.Printf("Build version: %s\n", validateBuildInfo(buildVersion))
	fmt.Printf("Build date: %s\n", validateBuildInfo(buildDate))
	fmt.Printf("Build commit: %s\n", validateBuildInfo(buildCommit))
}

func main() {
	printBuildInfo()
	cfg := config.NewConfig()
	var repo repository.RepositoryInterface
	var err error
	auditors := make([]service.Auditor, 0)
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
	if cfg.AuditFile != "" {
		auditor, err := service.NewInfileAuditor(cfg.AuditFile)
		if err != nil {
			panic(err)
		}
		auditors = append(auditors, auditor)
	}
	if cfg.AuditURL != "" {
		auditor := service.NewRemoteAuditor(cfg.AuditURL)
		auditors = append(auditors, auditor)
	}
	logger.Log.Info("Repository type set to", zap.String("type", cfg.RepositoryType))
	app := NewApp(repo, *cfg)
	appHandler := app.GetHandler(auditors)
	if cfg.EnableHTTPS {
		certFile := "cert/cert.pem"
		keyFile := "cert/private.pem"
		logger.Log.Info("Running HTTPS server", zap.String("address", cfg.ServerHostPort))
		err = http.ListenAndServeTLS(cfg.ServerHostPort, certFile, keyFile, logger.RequestLogger(handler.AuthMiddleware(handler.GzipMiddleware(appHandler), repo)))
	} else {
		logger.Log.Info("Running HTTP server", zap.String("address", cfg.ServerHostPort))
		err = http.ListenAndServe(cfg.ServerHostPort, logger.RequestLogger(handler.AuthMiddleware(handler.GzipMiddleware(appHandler), repo)))
	}
	if err != nil {
		logger.Log.Fatal("Server failed", zap.Error(err))
		panic(err)
	}
}
