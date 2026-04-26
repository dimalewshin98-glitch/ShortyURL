package repository

import (
	"context"
	"sync"
)

type UrlInfo struct {
	URL    string
	UserID string
}

type InmemoryRepository struct {
	urls map[string]UrlInfo
	mu   sync.RWMutex
}

func NewInmemoryRepository() *InmemoryRepository {
	return &InmemoryRepository{
		urls: make(map[string]UrlInfo),
	}
}

func (r *InmemoryRepository) Ping(ctx context.Context) error {
	return nil
}

func (r *InmemoryRepository) Store(ctx context.Context, userID string, urlID string, URL string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.urls[urlID] = UrlInfo{URL: URL, UserID: userID}
	return urlID, nil
}

func (r *InmemoryRepository) Get(ctx context.Context, urlID string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	URLInfo := r.urls[urlID]
	return URLInfo.URL, nil
}
