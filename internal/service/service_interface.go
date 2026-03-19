package service

type ServiceInterface interface {
	Shorten(url string) (string, error)
	GetUrl(urlId string) (string, error)
}
