package repository

import (
	"context"

	models "github.com/dimalewshin98-glitch/ShortyURL/internal/model"
)

type RepositoryInterface interface {
	Store(ctx context.Context, userID int, urlId string, url string) (string, error)
	Get(ctx context.Context, urlId string) (string, error)
	GetUserUrls(ctx context.Context, userID int) (models.ApiUserUrlsRes, error)
	Ping(ctx context.Context) error
	GetUsersID(ctx context.Context) ([]int, error)
}
