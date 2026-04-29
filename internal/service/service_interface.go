package service

import (
	"context"

	models "github.com/dimalewshin98-glitch/ShortyURL/internal/model"
)

type ServiceInterface interface {
	Shorten(ctx context.Context, userID int, url string) (string, error)
	ShortenBatch(ctx context.Context, userID int, url models.ApiShortenBatchReq) (models.ApiShortenBatchRes, error)
	GetURL(ctx context.Context, urlId string) (string, error)
	UserUrls(ctx context.Context, userID int) (models.ApiUserUrlsRes, error)
	Ping(ctx context.Context) error
}
