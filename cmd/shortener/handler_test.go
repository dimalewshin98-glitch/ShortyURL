package main

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dimalewshin98-glitch/ShortyURL/internal/handler"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/repository"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/service"
	"github.com/stretchr/testify/assert"
)

type MockedInmemoryRepository struct {
	*repository.InmemoryRepository
}

func (mir *MockedInmemoryRepository) Get(urlID string) (string, error) {
	return "https://mockedurl.com", nil
}

func (mir *MockedInmemoryRepository) Store(urlId string, url string) (string, error) {
	return "AbCdEf", nil
}

func TestStorenHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response    string
	}
	tests := []struct {
		name        string
		contentType string
		body        string
		request     string
		requestType string
		want        want
	}{
		{
			name:        "test 1 | Success",
			contentType: "text/plain",
			body:        "https://mockedurl.com",
			want: want{
				contentType: "text/plain",
				statusCode:  201,
				response:    "http://localhost:8080/AbCdEf",
			},
			request:     "/",
			requestType: "POST",
		},
		{
			name:        "test 2 | Unsuccess | Request type error",
			contentType: "text/plain",
			body:        "https://mockedurl.com",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				response:    "Method not allowed\n",
			},
			request:     "/",
			requestType: "GET",
		},
		{
			name:        "test 3 | Unsuccess | Content-Type error",
			contentType: "text/html",
			body:        "https://mockedurl.com",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				response:    "Content-Type not allowed\n",
			},
			request:     "/",
			requestType: "POST",
		},
		{
			name:        "test 4 | Unsuccess | URL is empty",
			contentType: "text/plain",
			body:        "",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				response:    "URL is empty\n",
			},
			request:     "/",
			requestType: "POST",
		},
	}
	repository := repository.NewInmemoryRepository()
	mockedRepository := &MockedInmemoryRepository{repository}
	shorterService := service.NewShorterService(mockedRepository)
	requestsHandler := handler.NewRequestsHandler(shorterService)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.requestType, tt.request, strings.NewReader(tt.body))
			request.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			requestsHandler.Shorten(w, request)
			resBytes, _ := io.ReadAll(w.Body)
			assert.Equal(t, tt.want.statusCode, w.Result().StatusCode)
			assert.Equal(t, tt.want.response, string(resBytes))
			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))
		})
	}
}

func TestGetUrlHandler(t *testing.T) {
	type want struct {
		statusCode int
		response   string
	}
	tests := []struct {
		name        string
		contentType string
		request     string
		requestType string
		want        want
	}{
		{
			name:        "test 1 | Success",
			contentType: "text/plain",
			want: want{
				statusCode: 307,
				response:   "Location: https://mockedurl.com",
			},
			request:     "/AbCdEf",
			requestType: "GET",
		},
		{
			name:        "test 2 | Unsuccess | Request type error",
			contentType: "text/plain",
			want: want{
				statusCode: 400,
				response:   "Method not allowed\n",
			},
			request:     "/AbCdEf",
			requestType: "POST",
		},
		{
			name:        "test 3 | Unsuccess | Content-Type error",
			contentType: "text/html",
			want: want{
				statusCode: 400,
				response:   "Content-Type not allowed\n",
			},
			request:     "/AbCdEf",
			requestType: "GET",
		},
	}
	repository := repository.NewInmemoryRepository()
	mockedRepository := &MockedInmemoryRepository{repository}
	shorterService := service.NewShorterService(mockedRepository)
	requestsHandler := handler.NewRequestsHandler(shorterService)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.requestType, tt.request, nil)
			request.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			requestsHandler.GetUrl(w, request)
			resBytes, _ := io.ReadAll(w.Body)
			assert.Equal(t, tt.want.statusCode, w.Result().StatusCode)
			assert.Equal(t, tt.want.response, string(resBytes))
		})
	}
}
