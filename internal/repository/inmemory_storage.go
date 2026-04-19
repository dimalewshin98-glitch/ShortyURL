package repository

import (
	"context"
	"sync"
)

type InmemoryRepository struct {
	urls map[string]string
	mu   sync.RWMutex
}

func NewInmemoryRepository() *InmemoryRepository {
	return &InmemoryRepository{
		urls: make(map[string]string),
	}
}

func (r *InmemoryRepository) Ping(ctx context.Context) error {
	return nil
}

func (r *InmemoryRepository) Store(ctx context.Context, urlID string, URL string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.urls[urlID] = URL
	return urlID, nil
}

func (r *InmemoryRepository) Get(ctx context.Context, urlID string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	URL := r.urls[urlID]
	return URL, nil
}
