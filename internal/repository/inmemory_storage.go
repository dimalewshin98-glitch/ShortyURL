package repository

import (
	"context"
	"sync"

	models "github.com/dimalewshin98-glitch/ShortyURL/internal/model"
)

type UrlInfo struct {
	URL    string
	UserID int
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

func (r *InmemoryRepository) Store(ctx context.Context, userID int, urlID string, URL string) (string, error) {
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

func (r *InmemoryRepository) GetUsersID(ctx context.Context) ([]int, error) {
	uniqueIDs := make(map[int]bool)
	var result []int
	for _, urlInfo := range r.urls {
		if !uniqueIDs[urlInfo.UserID] {
			uniqueIDs[urlInfo.UserID] = true
			result = append(result, urlInfo.UserID)
		}
	}
	return result, nil
}

func (r *InmemoryRepository) GetUserUrls(ctx context.Context, userID int) (models.ApiUserUrlsRes, error) {
	var userURLs models.ApiUserUrlsRes
	for shortURL, urlInfo := range r.urls {
		var userURL models.UserUrlRes
		if urlInfo.UserID == userID {
			userURL.OriginalURL = urlInfo.URL
			userURL.ShortURL = shortURL
			userURLs = append(userURLs, userURL)
		}
	}
	return userURLs, nil
}
