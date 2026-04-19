package repository

import "context"

type RepositoryInterface interface {
	Store(ctx context.Context, urlId string, url string) (string, error)
	Get(ctx context.Context, urlId string) (string, error)
	Ping(ctx context.Context) error
}
