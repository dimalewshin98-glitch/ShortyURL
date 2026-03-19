package repository

type InmemoryRepository struct {
	urls map[string]string
}

func NewInmemoryRepository() *InmemoryRepository {
	return &InmemoryRepository{
		urls: make(map[string]string),
	}
}

func (r *InmemoryRepository) Store(urlId string, url string) {
	r.urls[urlId] = url
}

func (r *InmemoryRepository) Get(urlId string) string {
	url := r.urls[urlId]
	return url
}
