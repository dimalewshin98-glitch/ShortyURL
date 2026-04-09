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

func (s *ShorterService) Ping() error {
	err := s.repo.Ping()
	return err
}

func (s *ShorterService) GetURL(urlID string) (string, error) {
	URL, err := s.repo.Get(urlID)
	if err != nil {
		return "", err
	}
	if URL == "" {
		return URL, errors.New("URL not found")
	}
	return URL, nil
}

func (s *ShorterService) Shorten(URL string) (string, error) {
	urlID, err := s.repo.Store(s.generateShortURL(), URL)
	if err != nil {
		return "", err
	}
	return s.config.ShortenURLHostPort + "/" + urlID, nil
}

func (s *ShorterService) generateShortURL() string {
	urlIDLen := 6
	letters := "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	result := make([]byte, urlIDLen)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
