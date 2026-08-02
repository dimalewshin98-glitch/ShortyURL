package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	models "github.com/dimalewshin98-glitch/ShortyURL/internal/model"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/repository"
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

// Ping — обработчик для проверки доступности сервиса.
// Обрабатывает только GET-запросы.
// В случае успеха возвращает HTTP-статус 200 OK.
// При ошибке в работе сервиса возвращает 500 Internal Server Error.
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

// ApiInternalStats — обработчик для получения внутренней статистики сервиса.
// Обрабатывает только GET-запросы.
// В случае успеха возвращает HTTP-статус 200 OK и JSON-объект со статистикой.
// При ошибке в работе сервиса возвращает 500 Internal Server Error.
func (s *RequestsHandler) ApiInternalStats(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}
	internalStats, err := s.service.InternalStats(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	err = enc.Encode(internalStats)
	if err != nil {
		http.Error(w, "Json response encode error", http.StatusInternalServerError)
		return
	}
}

// GetURL — обработчик для перенаправления пользователя по короткому URL.
// Обрабатывает только GET-запросы с Content-Type: text/plain.
// Извлекает ID короткого URL из пути и выполняет перенаправление (307 Temporary Redirect) на исходный адрес.
// Если URL был удален, возвращает 410 Gone.
func (s *RequestsHandler) GetURL(w http.ResponseWriter, r *http.Request) {
	rawUserID := r.Context().Value("userID")
	if rawUserID == nil {
		http.Error(w, "User ID not found in context", http.StatusBadRequest)
		return
	}
	userID, ok := rawUserID.(int)
	if !ok {
		http.Error(w, "Invalid User ID type in context", http.StatusBadRequest)
		return
	}
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
	URL, err := s.service.GetURL(ctx, userID, urlID)
	if err != nil {
		if errors.Is(err, repository.ErrShortURLDeleted) {
			w.WriteHeader(http.StatusGone)
			return
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	w.Header().Set("Location", URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// Shorten — обработчик для создания короткого URL из тела POST-запроса.
// Принимает на вход текстовый (text/plain) POST-запрос, где тело — это исходный URL.
// В случае успеха возвращает созданный короткий URL в теле ответа со статусом 201 Created.
// Если короткий URL уже существует, возвращает статус 409 Conflict.
func (s *RequestsHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	rawUserID := r.Context().Value("userID")
	if rawUserID == nil {
		http.Error(w, "User ID not found in context", http.StatusBadRequest)
		return
	}
	userID, ok := rawUserID.(int)
	if !ok {
		http.Error(w, "Invalid User ID type in context", http.StatusBadRequest)
		return
	}
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

// ApiShorten — JSON API-обработчик для создания одного короткого URL.
// Принимает POST-запрос с телом в формате models.ApiShortenReq.
// Возвращает ответ в формате models.ApiShortenRes со статусом 201 Created.
// При ошибке возвращает 400 Bad Request.
func (s *RequestsHandler) ApiShorten(w http.ResponseWriter, r *http.Request) {
	rawUserID := r.Context().Value("userID")
	if rawUserID == nil {
		http.Error(w, "User ID not found in context", http.StatusBadRequest)
		return
	}
	userID, ok := rawUserID.(int)
	if !ok {
		http.Error(w, "Invalid User ID type in context", http.StatusBadRequest)
		return
	}
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
		if errors.Is(err, repository.ErrShortURLExists) {
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

// ApiShortenBatch — JSON API-обработчик для пакетного создания коротких URL.
// Принимает POST-запрос с телом в формате models.ApiShortenBatchReq (список URL).
// Возвращает результат в формате models.ApiShortenBatchRes со статусом 201 Created.
// При ошибке возвращает 400 Bad Request.
func (s *RequestsHandler) ApiShortenBatch(w http.ResponseWriter, r *http.Request) {
	rawUserID := r.Context().Value("userID")
	if rawUserID == nil {
		http.Error(w, "User ID not found in context", http.StatusBadRequest)
		return
	}
	userID, ok := rawUserID.(int)
	if !ok {
		http.Error(w, "Invalid User ID type in context", http.StatusBadRequest)
		return
	}
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
	if err != nil && !errors.Is(err, repository.ErrShortURLExists) {
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

// ApiUserUrls — JSON API-обработчик для получения списка всех коротких URL пользователя.
// Обрабатывает GET-запросы.
// Возвращает список URL в формате models.ApiUserUrlsRes.
// При ошибке возвращает 400 Bad Request.
func (s *RequestsHandler) ApiUserUrls(w http.ResponseWriter, r *http.Request) {
	rawUserID := r.Context().Value("userID")
	if rawUserID == nil {
		http.Error(w, "User ID not found in context", http.StatusBadRequest)
		return
	}
	userID, ok := rawUserID.(int)
	if !ok {
		http.Error(w, "Invalid User ID type in context", http.StatusBadRequest)
		return
	}
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

// Delete — JSON API-обработчик для удаления одного или нескольких коротких URL.
// Принимает DELETE-запрос с телом в формате models.ApiDeleteReq, содержащим список ID для удаления.
// Возвращает ответ со статусом 202 Accepted.
// При ошибке возвращает 400 Bad Request.
func (s *RequestsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	rawUserID := r.Context().Value("userID")
	if rawUserID == nil {
		http.Error(w, "User ID not found in context", http.StatusBadRequest)
		return
	}
	userID, ok := rawUserID.(int)
	if !ok {
		http.Error(w, "Invalid User ID type in context", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type not allowed", http.StatusBadRequest)
		return
	}
	var req models.ApiDeleteReq
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := dec.Decode(&req)
	if err != nil {
		http.Error(w, "Json request decode error", http.StatusBadRequest)
		return
	}
	_, err = s.service.Delete(ctx, req, userID)
	resHeader := http.StatusAccepted
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resHeader)
}
