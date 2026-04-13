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
	DatabaseDsn        string
}

func NewConfig() *Config {
	serverHostPort := flag.String("a", "localhost:8888", "server host:port")
	shortenURLHostPort := flag.String("b", "http://localhost:8000", "shorten url http://host:port")
	logLevel := flag.String("l", "info", "log level")
	fileStoragePath := flag.String("f", "", "file storage path")
	databaseDsn := flag.String("d", "", "databse destination (host/host:port)")
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
	if envDatabaseDsn := os.Getenv("DATABASE_DSN"); envDatabaseDsn != "" {
		*databaseDsn = envDatabaseDsn
	}
	conf := &Config{
		ServerHostPort:     *serverHostPort,
		ShortenURLHostPort: *shortenURLHostPort,
		LogLevel:           *logLevel,
		FileStoragePath:    *fileStoragePath,
		DatabaseDsn:        *databaseDsn,
	}
	conf.RepositoryType = conf.setRepositoryType(*fileStoragePath, *databaseDsn)
	return conf
}

func (c *Config) setRepositoryType(fileStoragePath string, databaseDsn string) string {
	if databaseDsn != "" {
		return "db"
	}
	if fileStoragePath != "" {
		return "file"
	}
	return "memory"
}
