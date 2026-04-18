package service

import (
	"context"
	"errors"
	"math/rand"

	"github.com/dimalewshin98-glitch/ShortyURL/internal/config"
	models "github.com/dimalewshin98-glitch/ShortyURL/internal/model"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/repository"
	"github.com/google/uuid"
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

func (s *ShorterService) Ping(ctx context.Context) error {
	err := s.repo.Ping(ctx)
	return err
}

func (s *ShorterService) GetURL(ctx context.Context, urlID string) (string, error) {
	URL, err := s.repo.Get(ctx, urlID)
	if err != nil {
		return "", err
	}
	if URL == "" {
		return URL, errors.New("URL not found")
	}
	return URL, nil
}

func (s *ShorterService) Shorten(ctx context.Context, URL string) (string, error) {
	urlID, err := s.repo.Store(ctx, s.generateShortURL(), URL)
	if err != nil && err.Error() != "Short URL already exists" {
		return "", err
	}
	return s.config.ShortenURLHostPort + "/" + urlID, err
}

func (s *ShorterService) ShortenBatch(ctx context.Context, reqData models.ApiShortenBatchReq) (models.ApiShortenBatchRes, error) {
	var err error = nil
	ctxUUID := uuid.New().String()
	ctxVal := context.WithValue(ctx, "UUID", ctxUUID)
	resData := models.ApiShortenBatchRes{}
	for i := range reqData {
		if i == len(reqData)-1 {
			ctxVal = context.WithValue(ctxVal, "isLastReq", true)
		}
		urlID, err := s.repo.Store(ctxVal, s.generateShortURL(), reqData[i].OriginalURL)
		if err != nil && err.Error() != "Short URL already exists" {
			return models.ApiShortenBatchRes{}, err
		}
		shortURL := s.config.ShortenURLHostPort + "/" + urlID
		resData = append(resData, models.BatchItemRes{CorrelationID: reqData[i].CorrelationID, ShortURL: shortURL})
	}
	return resData, err
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
