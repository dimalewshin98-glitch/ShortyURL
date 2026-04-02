package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerHostPort     string
	ShortenURLHostPort string
	LogLevel           string
	FileStoragePath    string
	RepositoryType     string
}

func NewConfig() *Config {
	serverHostPort := flag.String("a", "localhost:8888", "server host:port")
	shortenURLHostPort := flag.String("b", "http://localhost:8000", "shorten url http://host:port")
	logLevel := flag.String("l", "info", "log level")
	fileStoragePath := flag.String("f", "storage", "file storage path")
	repositoryType := flag.String("r", "file", "repository type (file/memory)")
	flag.Parse()
	if envServerHostPort := os.Getenv("SERVER_ADDRESS"); envServerHostPort != "" {
		*serverHostPort = envServerHostPort
	}
	if envShortenURLHostPort := os.Getenv("BASE_URL"); envShortenURLHostPort != "" {
		*shortenURLHostPort = envShortenURLHostPort
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		*logLevel = envLogLevel
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		*fileStoragePath = envFileStoragePath
	}
	if envRepositoryType := os.Getenv("REPOSITORY TYPE"); envRepositoryType != "" {
		*repositoryType = envRepositoryType
	}
	return &Config{
		ServerHostPort:     *serverHostPort,
		ShortenURLHostPort: *shortenURLHostPort,
		LogLevel:           *logLevel,
		FileStoragePath:    *fileStoragePath,
		RepositoryType:     *repositoryType,
	}
}
