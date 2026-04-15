package repository

import "context"

type InmemoryRepository struct {
	urls map[string]string
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
	r.urls[urlID] = URL
	return urlID, nil
}

func (r *InmemoryRepository) Get(ctx context.Context, urlID string) (string, error) {
	URL := r.urls[urlID]
	return URL, nil
}
