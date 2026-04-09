package repository

type InmemoryRepository struct {
	urls map[string]string
}

func NewInmemoryRepository() *InmemoryRepository {
	return &InmemoryRepository{
		urls: make(map[string]string),
	}
}

func (r *InmemoryRepository) Ping() error {
	return nil
}

func (r *InmemoryRepository) Store(urlID string, URL string) (string, error) {
	r.urls[urlID] = URL
	return urlID, nil
}

func (r *InmemoryRepository) Get(urlID string) (string, error) {
	URL := r.urls[urlID]
	return URL, nil
}
