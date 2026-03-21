package service

type ServiceInterface interface {
	Shorten(url string) (string, error)
	GetURL(urlId string) (string, error)
}
