package repository

import (
	"context"
	"errors"

	models "github.com/dimalewshin98-glitch/ShortyURL/internal/model"
)

var ErrShortURLExists = errors.New("short URL already exists")
var ErrShortURLDeleted = errors.New("short URL already deleted")

type RepositoryInterface interface {
	Store(ctx context.Context, userID int, urlID string, url string) (string, error)
	Get(ctx context.Context, urlID string) (string, error)
	GetUserUrls(ctx context.Context, userID int) (models.ApiUserUrlsRes, error)
	Ping(ctx context.Context) error
	GetUsersID(ctx context.Context) ([]int, error)
	SetDelete(ctx context.Context, userID int, urlID string) (string, error)
}
