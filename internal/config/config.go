package config

import (
	"encoding/json"
	"flag"
	"os"
)

type Config struct {
	ServerHostPort     string
	EnableHTTPS        bool
	ShortenURLHostPort string
	LogLevel           string
	FileStoragePath    string
	RepositoryType     string
	DatabaseDsn        string
	AuditFile          string
	AuditURL           string
	TrustedSubnet      string
}

type JSONConfig struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
	LogLevel        string `json:"log_level"`
	AuditFile       string `json:"audit_file"`
	AuditURL        string `json:"audit_url"`
	TrustedSubnet   string `json:"trusted_subnet"`
}

func NewConfig() (*Config, error) {
	var configFile string
	flag.StringVar(&configFile, "c", "", "config file path")
	flag.StringVar(&configFile, "config", "", "config file path")
	serverHostPort := flag.String("a", "localhost:8888", "server host:port")
	enableHTTPS := flag.Bool("s", false, "enable HTTPS")
	shortenURLHostPort := flag.String("b", "http://localhost:8000", "shorten url http://host:port")
	logLevel := flag.String("l", "info", "log level")
	fileStoragePath := flag.String("f", "", "file storage path")
	databaseDsn := flag.String("d", "", "databse destination (host/host:port)")
	trustedSubnet := flag.String("t", "", "trusted subnet")
	auditFile := flag.String("audit-file", "", "audit file")
	auditURL := flag.String("audit-url", "", "audit url")
	flag.Parse()
	visitedFlags := make(map[string]bool)
	flag.VisitAll(func(f *flag.Flag) {
		visitedFlags[f.Name] = true
	})
	if envConfigFile := os.Getenv("CONFIG"); envConfigFile != "" {
		configFile = envConfigFile
	}
	var jsonConfig JSONConfig
	if configFile != "" {
		if err := loadJSONConfig(configFile, &jsonConfig); err != nil {
			return nil, err
		}
	}
	if !visitedFlags["a"] && jsonConfig.ServerAddress != "" {
		*serverHostPort = jsonConfig.ServerAddress
	}
	if !visitedFlags["b"] && jsonConfig.BaseURL != "" {
		*shortenURLHostPort = jsonConfig.BaseURL
	}
	if !visitedFlags["f"] && jsonConfig.FileStoragePath != "" {
		*fileStoragePath = jsonConfig.FileStoragePath
	}
	if !visitedFlags["f"] && jsonConfig.DatabaseDSN != "" {
		*databaseDsn = jsonConfig.DatabaseDSN
	}
	if !visitedFlags["s"] && jsonConfig.EnableHTTPS {
		*enableHTTPS = jsonConfig.EnableHTTPS
	}
	if !visitedFlags["l"] && jsonConfig.LogLevel != "" {
		*logLevel = jsonConfig.LogLevel
	}
	if !visitedFlags["t"] && jsonConfig.TrustedSubnet != "" {
		*trustedSubnet = jsonConfig.TrustedSubnet
	}
	if !visitedFlags["audit-file"] && jsonConfig.AuditFile != "" {
		*auditFile = jsonConfig.AuditFile
	}
	if !visitedFlags["audit-url"] && jsonConfig.AuditURL != "" {
		*auditURL = jsonConfig.AuditURL
	}
	if envServerHostPort := os.Getenv("SERVER_ADDRESS"); envServerHostPort != "" {
		*serverHostPort = envServerHostPort
	}
	if envEnableHTTPS := os.Getenv("ENABLE_HTTPS"); envEnableHTTPS != "" {
		*enableHTTPS = envEnableHTTPS == "true"
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
	if envAuditFile := os.Getenv("AUDIT_FILE"); envAuditFile != "" {
		*auditFile = envAuditFile
	}
	if envAuditURL := os.Getenv("AUDIT_URL"); envAuditURL != "" {
		*auditURL = envAuditURL
	}
	if envTrustedSubnet := os.Getenv("TRUSTED_SUBNET"); envTrustedSubnet != "" {
		*trustedSubnet = envTrustedSubnet
	}
	conf := &Config{
		ServerHostPort:     *serverHostPort,
		EnableHTTPS:        *enableHTTPS,
		ShortenURLHostPort: *shortenURLHostPort,
		LogLevel:           *logLevel,
		FileStoragePath:    *fileStoragePath,
		DatabaseDsn:        *databaseDsn,
		AuditFile:          *auditFile,
		AuditURL:           *auditURL,
		TrustedSubnet:      *trustedSubnet,
	}
	conf.RepositoryType = conf.setRepositoryType(*fileStoragePath, *databaseDsn)
	return conf, nil
}

func loadJSONConfig(filename string, config *JSONConfig) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, config); err != nil {
		return err
	}
	return nil
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
