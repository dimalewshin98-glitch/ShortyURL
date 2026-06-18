package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/dimalewshin98-glitch/ShortyURL/internal/config"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/handler"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/mocks"
	models "github.com/dimalewshin98-glitch/ShortyURL/internal/model"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/repository"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/service"
	"github.com/golang/mock/gomock"
)

// func (mir *MockedInmemoryRepository) Get(ctx context.Context, urlID string) (string, error) {
// 	return "https://mockedurl.com", nil
// }

// func (mir *MockedInmemoryRepository) Store(ctx context.Context, userID int, urlID string, URL string) (string, error) {
// 	return "AbCdEf", nil
// }

// func (mir *MockedInmemoryRepository) GetUserUrls(ctx context.Context, userID int) (models.ApiUserUrlsRes, error) {
// 	var res models.ApiUserUrlsRes
// 	res = append(res, models.UserUrlRes{ShortURL: "a", OriginalURL: "b"})
// 	res = append(res, models.UserUrlRes{ShortURL: "c", OriginalURL: "d"})
// 	return res, nil
// }

// func (mir *MockedInmemoryRepository) Ping(ctx context.Context) error {
// 	return nil
// }

// func (mir *MockedInmemoryRepository) GetUsersID(ctx context.Context) ([]int, error) {
// 	return []int{1, 2, 3}, nil
// }

func ExampleRequestsHandler_Ping() {
	type want struct {
		statusCode int
		response   string
		dbResponse error
	}
	test := struct {
		request     string
		requestType string
		want        want
	}{
		request:     "/ping",
		requestType: "GET",
		want: want{
			statusCode: 200,
			response:   "",
			dbResponse: nil,
		},
	}
	ctrl := gomock.NewController(nil)
	mockedRepository := mocks.NewMockRepositoryInterface(ctrl)
	if test.requestType == "GET" {
		mockedRepository.EXPECT().
			Ping(gomock.Any()).
			Return(test.want.dbResponse)
	}
	mockedConfig := &config.Config{
		ServerHostPort:     "localhost:8080",
		ShortenURLHostPort: "http://localhost:8080",
	}
	shorterService := service.NewShorterService(mockedRepository, mockedConfig)
	requestsHandler := handler.NewRequestsHandler(shorterService)
	request := httptest.NewRequest(test.requestType, test.request, nil)
	w := httptest.NewRecorder()
	requestsHandler.Ping(w, request)
	resBytes, _ := io.ReadAll(w.Body)
	fmt.Println(w.Code)
	if len(resBytes) > 0 {
		fmt.Println(string(resBytes))
	}

	// Output:
	// 200
}

func ExampleRequestsHandler_ApiShortenBatch() {
	type want struct {
		contentType     string
		statusCode      int
		response        string
		contentEncoding string
	}
	test := struct {
		contentType     string
		acceptEncoding  string
		contentEncoding string
		body            string
		request         string
		requestType     string
		want            want
	}{
		contentType: "application/json",
		body: `[
						{
							"correlation_id": "aaa",
							"original_url": "urlAA"
						},
						{
							"correlation_id": "aaa",
							"original_url": "urlAA"
						}
					] `,
		want: want{
			contentType: "application/json",
			statusCode:  201,
			response:    `[{"correlation_id":"aaa","short_url":"http://localhost:8080/bjjBrD"},{"correlation_id":"aaa","short_url":"http://localhost:8080/bjjBrD"}]` + "\n",
		},
		request:     "/api/shorten/batch",
		requestType: "POST",
	}

	ctrl := gomock.NewController(nil)
	mockedRepository := mocks.NewMockRepositoryInterface(ctrl)
	mockedRepository.EXPECT().
		GetUsersID(gomock.Any()).
		Return([]int{1, 2, 3}, nil).
		AnyTimes()
	if test.requestType == "POST" && test.want.statusCode != 400 {
		mockedRepository.EXPECT().
			Store(gomock.Any(), 4, gomock.Any(), "urlAA").
			Return("bjjBrD", nil).
			Times(2)
	}
	mockedConfig := &config.Config{
		ServerHostPort:     "localhost:8080",
		ShortenURLHostPort: "http://localhost:8080",
	}
	shorterService := service.NewShorterService(mockedRepository, mockedConfig)
	requestsHandler := handler.NewRequestsHandler(shorterService)
	appHandler := http.HandlerFunc(requestsHandler.ApiShortenBatch)
	handlerWithMiddleware := handler.AuthMiddleware(appHandler, mockedRepository)
	request := httptest.NewRequest(test.requestType, test.request, strings.NewReader(test.body))
	request.Header.Set("content-Type", test.contentType)
	request.Header.Set("Accept-Encoding", test.acceptEncoding)
	request.Header.Set("Content-Encoding", test.contentEncoding)
	w := httptest.NewRecorder()
	handlerWithMiddleware.ServeHTTP(w, request)
	resBytes, _ := io.ReadAll(w.Body)
	fmt.Println(w.Code)
	if len(resBytes) > 0 {
		fmt.Println(string(resBytes))
	}

	// Output:
	// 201
	// [{"correlation_id":"aaa","short_url":"http://localhost:8080/bjjBrD"},{"correlation_id":"aaa","short_url":"http://localhost:8080/bjjBrD"}]
}

func ExampleRequestsHandler_Shorten() {
	type want struct {
		contentType string
		statusCode  int
		response    string
	}
	test := struct {
		contentType string
		body        string
		request     string
		requestType string
		want        want
	}{
		contentType: "text/plain",
		body:        "urlAA",
		want: want{
			contentType: "text/plain",
			statusCode:  201,
			response:    "http://localhost:8080/bjjBrD",
		},
		request:     "/",
		requestType: "POST",
	}

	ctrl := gomock.NewController(nil)
	mockedRepository := mocks.NewMockRepositoryInterface(ctrl)
	mockedRepository.EXPECT().
		GetUsersID(gomock.Any()).
		Return([]int{1, 2, 3}, nil).
		AnyTimes()
	if test.requestType == "POST" && test.want.statusCode != 400 {
		mockedRepository.EXPECT().
			Store(gomock.Any(), 4, gomock.Any(), "urlAA").
			Return("bjjBrD", nil).
			Times(1)
	}
	mockedConfig := &config.Config{
		ServerHostPort:     "localhost:8080",
		ShortenURLHostPort: "http://localhost:8080",
	}
	shorterService := service.NewShorterService(mockedRepository, mockedConfig)
	requestsHandler := handler.NewRequestsHandler(shorterService)
	appHandler := http.HandlerFunc(requestsHandler.Shorten)
	handlerWithMiddleware := handler.AuthMiddleware(appHandler, mockedRepository)
	request := httptest.NewRequest(test.requestType, test.request, strings.NewReader(test.body))
	request.Header.Set("content-Type", test.contentType)
	w := httptest.NewRecorder()
	handlerWithMiddleware.ServeHTTP(w, request)
	resBytes, _ := io.ReadAll(w.Body)
	fmt.Println(w.Code)
	if len(resBytes) > 0 {
		fmt.Println(string(resBytes))
	}

	// Output:
	// 201
	// http://localhost:8080/bjjBrD
}

func ExampleRequestsHandler_ApiUserUrls() {
	type want struct {
		statusCode int
		response   string
	}

	test := struct {
		request     string
		requestType string
		want        want
	}{
		request:     "/api/user/urls",
		requestType: "GET",
		want: want{
			statusCode: 200,
			response:   `[{"short_url":"a","original_url":"b"},{"short_url":"c","original_url":"d"}]` + "\n",
		},
	}
	ctrl := gomock.NewController(nil)
	mockedRepository := mocks.NewMockRepositoryInterface(ctrl)
	mockedRepository.EXPECT().
		GetUsersID(gomock.Any()).
		Return([]int{1, 2, 3}, nil).
		AnyTimes()
	if test.requestType == "GET" && test.want.statusCode == 200 {
		var res models.ApiUserUrlsRes
		res = append(res, models.UserUrlRes{ShortURL: "a", OriginalURL: "b"})
		res = append(res, models.UserUrlRes{ShortURL: "c", OriginalURL: "d"})

		mockedRepository.EXPECT().
			GetUserUrls(gomock.Any(), gomock.Any()).
			Return(res, nil).
			Times(1)
	}
	mockedConfig := &config.Config{
		ServerHostPort:     "localhost:8080",
		ShortenURLHostPort: "http://localhost:8080",
	}
	shorterService := service.NewShorterService(mockedRepository, mockedConfig)
	requestsHandler := handler.NewRequestsHandler(shorterService)
	appHandler := http.HandlerFunc(requestsHandler.ApiUserUrls)
	handlerWithMiddleware := handler.AuthMiddleware(appHandler, mockedRepository)
	request := httptest.NewRequest(test.requestType, test.request, nil)
	w := httptest.NewRecorder()
	handlerWithMiddleware.ServeHTTP(w, request)
	resBytes, _ := io.ReadAll(w.Body)
	fmt.Println(w.Code)
	if len(resBytes) > 0 {
		fmt.Println(string(resBytes))
	}

	// Output:
	// 200
	// [{"short_url":"a","original_url":"b"},{"short_url":"c","original_url":"d"}]
}

func ExampleRequestsHandler_Delete() {
	type want struct {
		contentType     string
		statusCode      int
		response        string
		contentEncoding string
	}
	test := struct {
		contentType string
		body        string
		request     string
		requestType string
		want        want
	}{
		contentType: "application/json",
		body:        `["aaa"]`,
		request:     "/api/user/urls",
		requestType: "DELETE",
		want: want{
			contentType: "application/json",
			statusCode:  202,
			response:    ``,
		},
	}
	ctrl := gomock.NewController(nil)
	mockedRepository := mocks.NewMockRepositoryInterface(ctrl)
	mockedRepository.EXPECT().
		GetUsersID(gomock.Any()).
		Return([]int{1, 2, 3}, nil).
		AnyTimes()
	mockedRepository.EXPECT().
		SetDelete(gomock.Any(), 3, "aaa").
		Return("a", nil).
		Times(1)

	mockedConfig := &config.Config{
		ServerHostPort:     "localhost:8080",
		ShortenURLHostPort: "http://localhost:8080",
	}
	shorterService := service.NewShorterService(mockedRepository, mockedConfig)
	requestsHandler := handler.NewRequestsHandler(shorterService)
	appHandler := http.HandlerFunc(requestsHandler.Delete)
	handlerWithMiddleware := handler.AuthMiddleware(appHandler, mockedRepository)
	request := httptest.NewRequest(test.requestType, test.request, strings.NewReader(test.body))
	request.Header.Set("Content-Type", test.contentType)
	w := httptest.NewRecorder()
	handlerWithMiddleware.ServeHTTP(w, request)
	resBytes, _ := io.ReadAll(w.Body)
	fmt.Println(w.Code)
	if len(resBytes) > 0 {
		fmt.Println(string(resBytes))
	}

	// Output:
	// 202
}

func ExampleRequestsHandler_ApiShorten() {
	type want struct {
		contentType     string
		statusCode      int
		response        string
		contentEncoding string
	}
	test := struct {
		contentType     string
		body            string
		request         string
		requestType     string
		acceptEncoding  string
		contentEncoding string
		want            want
	}{
		contentType: "application/json",
		body:        `{"url": "https://mockedurl.com"}`,
		request:     "/api/shorten",
		requestType: "POST",
		want: want{
			contentType: "application/json",
			statusCode:  201,
			response:    `{"result":"http://localhost:8080/AbCdEf"}` + "\n",
		},
	}
	repository := repository.NewInmemoryRepository()
	mockedRepository := &MockedInmemoryRepository{repository}
	mockedConfig := &config.Config{
		ServerHostPort:     "localhost:8080",
		ShortenURLHostPort: "http://localhost:8080",
	}
	app := NewApp(mockedRepository, *mockedConfig)
	auditors := make([]service.Auditor, 0)
	appHandler := app.GetHandler(auditors)
	handlerWithMiddleware := handler.GzipMiddleware(handler.AuthMiddleware(appHandler, repository))
	var buf *bytes.Buffer
	if test.contentEncoding == "gzip" {
		var b bytes.Buffer
		zb := gzip.NewWriter(&b)
		zb.Write([]byte(test.body))
		zb.Close()
		buf = &b
	} else {
		buf = bytes.NewBufferString(test.body)
	}
	request := httptest.NewRequest(test.requestType, test.request, buf)
	request.Header.Set("Content-Type", test.contentType)
	request.Header.Set("Accept-Encoding", test.acceptEncoding)
	request.Header.Set("Content-Encoding", test.contentEncoding)
	w := httptest.NewRecorder()
	handlerWithMiddleware.ServeHTTP(w, request)
	var resBytes []byte
	if test.acceptEncoding == "gzip" {
		zr, _ := gzip.NewReader(w.Body)
		resBytes, _ = io.ReadAll(zr)
	} else {
		resBytes, _ = io.ReadAll(w.Body)
	}

	fmt.Println(w.Code)
	if len(resBytes) > 0 {
		fmt.Println(string(resBytes))
	}

	// Output:
	// 201
	// {"result":"http://localhost:8080/AbCdEf"}
}
