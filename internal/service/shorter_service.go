package service

import (
	"errors"
	"math/rand"

	"github.com/dimalewshin98-glitch/ShortyURL/internal/repository"
)

type ShorterService struct {
	repo *repository.InmemoryRepository
}

func NewShorterService(repo *repository.InmemoryRepository) *ShorterService {
	return &ShorterService{
		repo: repo,
	}
}

func (s *ShorterService) GetUrl(urlId string) (string, error) {
	url := s.repo.Get(urlId)
	if url == "" {
		return url, errors.New("URL not found")
	}
	return url, nil
}

func (s *ShorterService) Shorten(url string) (string, error) {
	urlId := s.generateShortUrl()
	s.repo.Store(urlId, url)
	return urlId, nil
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
