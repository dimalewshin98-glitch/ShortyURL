package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}
	err := s.service.Ping(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *RequestsHandler) GetURL(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Content-Type not allowed", http.StatusBadRequest)
		return
	}
	urlID := r.PathValue("id")
	URL, err := s.service.GetURL(ctx, urlID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *RequestsHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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
	shortURL, err := s.service.Shorten(ctx, userID, reqString)
	resHeader := http.StatusCreated
	if err != nil {
		if err.Error() == "Short URL already exists" {
			resHeader = http.StatusConflict
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(resHeader)
	w.Write([]byte(shortURL))
}

func (s *RequestsHandler) ApiShorten(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type not allowed", http.StatusBadRequest)
		return
	}
	var req models.ApiShortenReq
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
	shortURL, err := s.service.Shorten(ctx, userID, URL)
	resHeader := http.StatusCreated
	if err != nil {
		if err.Error() == "Short URL already exists" {
			resHeader = http.StatusConflict
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	res := models.ApiShortenRes{
		Result: shortURL,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resHeader)
	enc := json.NewEncoder(w)
	err = enc.Encode(res)
	if err != nil {
		http.Error(w, "Json response encode error", http.StatusBadRequest)
		return
	}
}

func (s *RequestsHandler) ApiShortenBatch(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type not allowed", http.StatusBadRequest)
		return
	}
	var req models.ApiShortenBatchReq
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := dec.Decode(&req)
	if err != nil {
		http.Error(w, "Json request decode error", http.StatusBadRequest)
		return
	}
	if len(req) == 0 {
		http.Error(w, "Batch is empty", http.StatusBadRequest)
		return
	}
	res, err := s.service.ShortenBatch(ctx, userID, req)
	if err != nil && err.Error() != "Short URL already exists" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
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

func (s *RequestsHandler) ApiUserUrls(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}
	res, err := s.service.UserUrls(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if res == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	err = enc.Encode(res)
	if err != nil {
		http.Error(w, "Json response encode error", http.StatusOK)
		return
	}
}
