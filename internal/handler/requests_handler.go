package handler

import (
	"encoding/json"
	"io"
	"net/http"

	models "github.com/dimalewshin98-glitch/ShortyURL/internal/model"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/service"
)

type RequestsHandler struct {
	service service.ServiceInterface
}

func NewRequestsHandler(service service.ServiceInterface) *RequestsHandler {
	return &RequestsHandler{
		service: service,
	}
}

func (s *RequestsHandler) Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}
	err := s.service.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *RequestsHandler) GetURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Content-Type not allowed", http.StatusBadRequest)
		return
	}
	urlID := r.PathValue("id")
	URL, err := s.service.GetURL(urlID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *RequestsHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Content-Type not allowed", http.StatusBadRequest)
		return
	}
	reqBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Reading body error", http.StatusBadRequest)
		return
	}
	reqString := string(reqBytes)
	if reqString == "" {
		http.Error(w, "URL is empty", http.StatusBadRequest)
		return
	}
	shortURL, err := s.service.Shorten(reqString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func (s *RequestsHandler) ApiShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type not allowed", http.StatusBadRequest)
		return
	}
	var req models.ApiRequest
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := dec.Decode(&req)
	if err != nil {
		http.Error(w, "Json request decode error", http.StatusBadRequest)
		return
	}
	URL := req.URL
	if URL == "" {
		http.Error(w, "URL is empty", http.StatusBadRequest)
		return
	}
	shortURL, err := s.service.Shorten(URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res := models.ApiResponse{
		Result: shortURL,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(w)
	err = enc.Encode(res)
	if err != nil {
		http.Error(w, "Json response encode error", http.StatusBadRequest)
		return
	}
}
