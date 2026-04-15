package service

import "context"

type ServiceInterface interface {
	Shorten(ctx context.Context, url string) (string, error)
	GetURL(ctx context.Context, urlId string) (string, error)
	Ping(ctx context.Context) error
}
