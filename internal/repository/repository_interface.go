package repository

type RepositoryInterface interface {
	Store(urlId string, url string) (string, error)
	Get(urlId string) (string, error)
}
