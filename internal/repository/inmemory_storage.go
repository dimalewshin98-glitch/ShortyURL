package repository

type InmemoryRepository struct {
	urls map[string]string
}

func NewInmemoryRepository() *InmemoryRepository {
	return &InmemoryRepository{
		urls: make(map[string]string),
	}
}

func (r *InmemoryRepository) Store(urlId string, url string) (string, error) {
	r.urls[urlId] = url
	return urlId, nil
}

func (r *InmemoryRepository) Get(urlId string) (string, error) {
	url := r.urls[urlId]
	return url, nil
}
