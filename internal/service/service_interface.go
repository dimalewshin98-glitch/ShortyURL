package service

type Service interface {
	Shorten(url string) (string, error)
	GetUrl(urlId string) (string, error)
}
