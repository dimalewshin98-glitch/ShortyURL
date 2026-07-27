package repository

import (
	"context"
	"errors"

	models "github.com/dimalewshin98-glitch/ShortyURL/internal/model"
)

var ErrShortURLExists = errors.New("short URL already exists")
var ErrShortURLDeleted = errors.New("short URL already deleted")

// RepositoryInterface описывает методы для взаимодействия с хранилищем данных.
// Реализация этого интерфейса отвечает за все операции CRUD (создание, чтение, обновление, удаление)
// с данными о коротких URL и их владельцах.
type RepositoryInterface interface {

	// Store сохраняет новую пару "короткий URL - оригинальный URL" для пользователя.
	//
	// Параметры:
	//   ctx - контекст запроса.
	//   ctxUUID - uuid контекста запроса.
	//   isLastReq - флаг последнего запроса
	//   userID - идентификатор пользователя, которому принадлежит URL.
	//   urlID - сгенерированный идентификатор короткой ссылки.
	//   url - оригинальный, полный URL-адрес.
	//
	// Возвращает:
	//   string - сохранённый короткий URL.
	//   error - ошибка, если URL уже существует или при записи в хранилище возникла проблема.
	Store(ctx context.Context, ctxUUID string, isLastReq bool, userID int, urlID string, url string) (string, error)

	// Get возвращает оригинальный URL по его короткому идентификатору (urlID).
	//
	// Параметры:
	//   ctx - контекст запроса.
	//   urlID - идентификатор короткой ссылки.
	//
	// Возвращает:
	//   string - оригинальный URL.
	//   error - ошибка, если запись с таким urlID не найдена.
	Get(ctx context.Context, urlID string) (string, error)

	// GetUserUrls возвращает список всех коротких URL-адресов, созданных указанным пользователем.
	//
	// Параметры:
	//   ctx - контекст запроса.
	//   userID - идентификатор пользователя.
	//
	// Возвращает:
	//   models.ApiUserUrlsRes - список структур UserUrlRes с парами "короткий-оригинальный" URL.
	//   error - ошибка при получении данных из хранилища.
	GetUserUrls(ctx context.Context, userID int) (models.ApiUserUrlsRes, error)

	// Ping проверяет доступность и работоспособность хранилища данных (например, соединение с БД).
	//
	// Параметры:
	//   ctx - контекст запроса.
	Ping(ctx context.Context) error

	// Close закрытие соединения с хранилицем данных (например, соединение с БД).
	//
	// Параметры:
	//   ctx - контекст запроса.
	Close(ctx context.Context) error

	// GetUsersID возвращает список всех идентификаторов пользователей, которые создавали короткие ссылки.
	//
	// Параметры:
	//   ctx - контекст запроса.
	GetUsersID(ctx context.Context) ([]int, error)

	// SetDelete помечает короткую ссылку как удалённую. Это "мягкое" удаление,
	// которое сохраняет данные в системе, но делает их недоступными для использования.
	//
	// Параметры:
	//   ctx - контекст запроса.
	//   ctxUUID - uuid контекста запроса.
	//   isLastReq - флаг последнего запроса
	//   userID - идентификатор пользователя-владельца ссылки.
	//   urlID - идентификатор короткой ссылки для удаления.
	SetDelete(ctx context.Context, ctxUUID string, isLastReq bool, userID int, urlID string) (string, error)
}
