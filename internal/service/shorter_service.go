package service

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/dimalewshin98-glitch/ShortyURL/internal/config"
	models "github.com/dimalewshin98-glitch/ShortyURL/internal/model"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/repository"
	"github.com/google/uuid"
)

type ShorterService struct {
	repo        repository.RepositoryInterface
	config      *config.Config
	msgChan     chan models.RepoDeleteMessage
	observers   []Auditor
	muObservers sync.RWMutex
}

func NewShorterService(repo repository.RepositoryInterface, config *config.Config) *ShorterService {
	serviceInstance := &ShorterService{
		repo:      repo,
		config:    config,
		msgChan:   make(chan models.RepoDeleteMessage, 1024),
		observers: make([]Auditor, 0),
	}
	go serviceInstance.flushMessages()
	return serviceInstance
}

func (s *ShorterService) Ping(ctx context.Context) error {
	err := s.repo.Ping(ctx)
	return err
}

func (s *ShorterService) GetURL(ctx context.Context, userID int, urlID string) (string, error) {
	URL, err := s.repo.Get(ctx, urlID)
	if err != nil {
		return "", err
	}
	if URL == "" {
		return URL, errors.New("URL not found")
	}
	go s.auditRequest("follow", userID, URL)
	return URL, nil
}

func (s *ShorterService) Shorten(ctx context.Context, userID int, URL string) (string, error) {
	urlID, err := s.repo.Store(ctx, userID, s.generateShortURL(), URL)
	if err != nil && !errors.Is(err, repository.ErrShortURLExists) {
		return "", err
	}
	go s.auditRequest("shorten", userID, URL)
	return s.config.ShortenURLHostPort + "/" + urlID, err
}

func (s *ShorterService) ShortenBatch(ctx context.Context, userID int, reqData models.ApiShortenBatchReq) (models.ApiShortenBatchRes, error) {
	var err error = nil
	ctxUUID := uuid.New().String()
	ctxVal := context.WithValue(ctx, "UUID", ctxUUID)
	resData := models.ApiShortenBatchRes{}
	for i := range reqData {
		if i == len(reqData)-1 {
			ctxVal = context.WithValue(ctxVal, "isLastReq", true)
		}
		urlID, err := s.repo.Store(ctxVal, userID, s.generateShortURL(), reqData[i].OriginalURL)
		if err != nil && !errors.Is(err, repository.ErrShortURLExists) {
			return models.ApiShortenBatchRes{}, err
		}
		shortURL := s.config.ShortenURLHostPort + "/" + urlID
		resData = append(resData, models.BatchItemRes{CorrelationID: reqData[i].CorrelationID, ShortURL: shortURL})
	}
	return resData, err
}

func (s *ShorterService) Delete(ctx context.Context, req models.ApiDeleteReq, userID int) ([]string, error) {
	var err error = nil
	resData := models.ApiDeleteRes{}
	for i := range req {
		s.msgChan <- models.RepoDeleteMessage{ShortURL: req[i], UserID: userID}
	}
	return resData, err
}

func (s *ShorterService) UserUrls(ctx context.Context, userID int) (models.ApiUserUrlsRes, error) {
	URL, err := s.repo.GetUserUrls(ctx, userID)
	if err != nil {
		return nil, err
	}
	return URL, nil
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

func (s *ShorterService) AttachAuditor(auditor Auditor) {
	s.muObservers.Lock()
	defer s.muObservers.Unlock()
	s.observers = append(s.observers, auditor)
}

func (s *ShorterService) auditRequest(enevtType string, iserID int, URL string) {
	s.muObservers.RLock()
	defer s.muObservers.RUnlock()
	for _, obs := range s.observers {
		obs.OnEvent(enevtType, iserID, URL)
	}
}

func (s *ShorterService) flushMessages() {
	ticker := time.NewTicker(10 * time.Second)
	var messages []models.RepoDeleteMessage
	for {
		select {
		case msg := <-s.msgChan:
			messages = append(messages, msg)
		case <-ticker.C:
			ctxUUID := uuid.New().String()
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			ctxVal := context.WithValue(ctx, "UUID", ctxUUID)
			for i := range messages {
				if i == len(messages)-1 {
					ctxVal = context.WithValue(ctxVal, "isLastReq", true)
				}
				_, err := s.repo.SetDelete(ctxVal, messages[i].UserID, messages[i].ShortURL)
				if err != nil {
					continue
				}
			}
			messages = nil
			cancel()
		}
	}
}
