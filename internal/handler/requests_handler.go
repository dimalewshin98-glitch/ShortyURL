package handler

import (
	"io"
	"net/http"

	"github.com/dimalewshin98-glitch/ShortyURL/internal/service"
)

type RequestsHandler struct {
	// service *service.ShorterService
	service service.ServiceInterface
}

// func NewRequestsHandler(service *service.ShorterService) *RequestsHandler {
func NewRequestsHandler(service service.ServiceInterface) *RequestsHandler {
	return &RequestsHandler{
		service: service,
	}
}

func (s *RequestsHandler) GetUrl(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Method not allowed", http.StatusBadRequest)
		return
	}
	if req.Header.Get("Content-Type") != "text/plain" {
		http.Error(res, "Content-Type not allowed", http.StatusBadRequest)
		return
	}
	urlId := req.PathValue("id")
	url, err := s.service.GetUrl(urlId)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusTemporaryRedirect)
	res.Write([]byte("Location: " + url))
}

func (s *RequestsHandler) Shorten(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method not allowed", http.StatusBadRequest)
		return
	}
	if req.Header.Get("Content-Type") != "text/plain" {
		http.Error(res, "Content-Type not allowed", http.StatusBadRequest)
		return
	}
	reqBytes, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(res, "Reading body error", http.StatusBadRequest)
		return
	}
	reqString := string(reqBytes)
	if reqString == "" {
		http.Error(res, "URL is empty", http.StatusBadRequest)
		return
	}
	shortUrl, err := s.service.Shorten(reqString)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortUrl))
}
