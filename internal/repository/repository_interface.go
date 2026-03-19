package repository

type Repository interface {
	Store(urlId string, url string)
	Get(urlId string) string
}
