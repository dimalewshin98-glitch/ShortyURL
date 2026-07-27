package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	var repo repository.RepositoryInterface
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
	app := NewApp(repo, *cfg)
	appHandler := app.GetHandler(auditors)
	var srv = http.Server{Addr: cfg.ServerHostPort, Handler: logger.RequestLogger(handler.AuthMiddleware(handler.GzipMiddleware(appHandler), repo))}
	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Log.Error("HTTP server Shutdown:", zap.Error(err))
		}
		if err := repo.Close(ctx); err != nil {
			logger.Log.Error("Repo Shutdown:", zap.Error(err))
		}
		close(idleConnsClosed)
	}()
	if cfg.EnableHTTPS {
		certFile := "cert/cert.pem"
		keyFile := "cert/private.pem"
		logger.Log.Info("Running HTTPS server", zap.String("address", cfg.ServerHostPort))
		if err := srv.ListenAndServeTLS(certFile, keyFile); err != http.ErrServerClosed {
			logger.Log.Error("Server failed", zap.Error(err))
			panic(err)
		}
	} else {
		logger.Log.Info("Running HTTP server", zap.String("address", cfg.ServerHostPort))
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Log.Error("Server failed", zap.Error(err))
			panic(err)
		}
	}
	<-idleConnsClosed
	logger.Log.Info("Server Shutdown gracefully", zap.String("address", cfg.ServerHostPort))
}
