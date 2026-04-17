package service

import (
	"context"

	models "github.com/dimalewshin98-glitch/ShortyURL/internal/model"
)

type ServiceInterface interface {
	Shorten(ctx context.Context, url string) (string, error)
	ShortenBatch(ctx context.Context, url models.ApiShortenBatchReq) (models.ApiShortenBatchRes, error)
	GetURL(ctx context.Context, urlId string) (string, error)
	Ping(ctx context.Context) error
}
