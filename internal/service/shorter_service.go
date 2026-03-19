package service

import (
	"errors"
	"math/rand"

	"github.com/dimalewshin98-glitch/ShortyURL/internal/config"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/repository"
)

type ShorterService struct {
	repo   repository.RepositoryInterface
	config *config.Config
}

func NewShorterService(repo repository.RepositoryInterface, config *config.Config) *ShorterService {
	return &ShorterService{
		repo:   repo,
		config: config,
	}
}

func (s *ShorterService) GetUrl(urlId string) (string, error) {
	url, err := s.repo.Get(urlId)
	if err != nil {
		return "", err
	}
	if url == "" {
		return url, errors.New("URL not found")
	}
	return url, nil
}

func (s *ShorterService) Shorten(url string) (string, error) {
	urlId, err := s.repo.Store(s.generateShortUrl(), url)
	if err != nil {
		return "", err
	}
	return s.config.ShortenUrlHostPort + "/" + urlId, nil
}

func (s *ShorterService) generateShortUrl() string {
	urlIdLen := 6
	letters := "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	result := make([]byte, urlIdLen)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
