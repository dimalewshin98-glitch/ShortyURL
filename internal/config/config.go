package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerHostPort     string
	ShortenURLHostPort string
}

func NewConfig() *Config {
	serverHostPort := flag.String("a", "localhost:8888", "server host:port")
	shortenURLHostPort := flag.String("b", "http://localhost:8000", "shorten url http://host:port")
	flag.Parse()
	if envServerHostPort := os.Getenv("SERVER_ADDRESS"); envServerHostPort != "" {
		*serverHostPort = envServerHostPort
	}
	if envShortenURLHostPort := os.Getenv("BASE_URL"); envShortenURLHostPort != "" {
		*shortenURLHostPort = envShortenURLHostPort
	}
	return &Config{
		ServerHostPort:     *serverHostPort,
		ShortenURLHostPort: *shortenURLHostPort,
	}
}
