package repository

import (
	"context"
	"sync"

	models "github.com/dimalewshin98-glitch/ShortyURL/internal/model"
)

type UrlInfo struct {
	URL       string
	UserID    int
	IsDeleted bool
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
	r.urls[urlID] = UrlInfo{URL: URL, UserID: userID, IsDeleted: false}
	return urlID, nil
}

func (r *InmemoryRepository) SetDelete(ctx context.Context, userID int, urlID string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	URLInfo := r.urls[urlID]
	if URLInfo.UserID == userID {
		URLInfo.IsDeleted = true
		r.urls[urlID] = URLInfo
	}
	return URLInfo.URL, nil
}

func (r *InmemoryRepository) Get(ctx context.Context, urlID string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	URLInfo := r.urls[urlID]
	isDeleted := URLInfo.IsDeleted
	if isDeleted {
		return URLInfo.URL, ErrShortURLDeleted
	}
	return URLInfo.URL, nil
}

func (r *InmemoryRepository) GetUsersID(ctx context.Context) ([]int, error) {
	uniqueIDs := make(map[int]struct{})
	for _, urlInfo := range r.urls {
		uniqueIDs[urlInfo.UserID] = struct{}{}
	}
	result := make([]int, 0, len(uniqueIDs))
	for userID := range uniqueIDs {
		result = append(result, userID)
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
