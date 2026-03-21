package config

import (
	"flag"
)

type Config struct {
	ServerHostPort     string
	ShortenURLHostPort string
}

func NewConfig() *Config {
	var serverHostPort = flag.String("a", "localhost:8888", "server host:port")
	var shortenURLHostPort = flag.String("b", "http://localhost:8000", "shorten url http://host:port")
	flag.Parse()
	return &Config{
		ServerHostPort:     *serverHostPort,
		ShortenURLHostPort: *shortenURLHostPort,
	}
}
