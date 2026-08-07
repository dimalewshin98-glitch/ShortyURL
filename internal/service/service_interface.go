package service

import (
	"context"

	models "github.com/dimalewshin98-glitch/ShortyURL/internal/model"
)

// ServiceInterface описывает бизнес-логику сервиса для управления короткими URL-адресами.
// Этот интерфейс используется HTTP-обработчиками (RequestsHandler) для выполнения операций.
type ServiceInterface interface {
	// Shorten создаёт одну короткую ссылку для указанного пользователя.
	//
	// Параметры:
	//   ctx - контекст запроса для управления таймаутами и отменой.
	//   userID - идентификатор пользователя, создающего ссылку.
	//   url - исходный длинный URL-адрес.
	//
	// Возвращает:
	//   string - сгенерированный короткий URL (ID).
	//   error - ошибка, если ссылка не может быть создана (например, уже существует).
	Shorten(ctx context.Context, userID int, url string) (string, error)

	// ShortenBatch создаёт несколько коротких ссылок в рамках одного запроса.
	//
	// Параметры:
	//   ctx - контекст запроса.
	//   userID - идентификатор пользователя.
	//   url - список моделей BatchItemReq, содержащих исходные URL и их correlation_id.
	//
	// Возвращает:
	//   models.ApiShortenBatchRes - список результатов с короткими URL и их correlation_id.
	//   error - ошибка, если пакетная операция не удалась.
	ShortenBatch(ctx context.Context, userID int, url models.ApiShortenBatchReq) (models.ApiShortenBatchRes, error)

	// GetURL возвращает оригинальный URL по его короткому идентификатору.
	//
	// Параметры:
	//   ctx - контекст запроса.
	//   userID - идентификатор пользователя, для которого выполняется поиск.
	//   urlId - идентификатор короткой ссылки (короткий URL).
	//
	// Возвращает:
	//   string - оригинальный длинный URL.
	//   error - ошибка, если URL не найден или не принадлежит пользователю.
	GetURL(ctx context.Context, userID int, urlId string) (string, error)

	// UserUrls возвращает список всех коротких URL-адресов, созданных пользователем.
	//
	// Параметры:
	//   ctx - контекст запроса.
	//   userID - идентификатор пользователя.
	//
	// Возвращает:
	//   models.ApiUserUrlsRes - список структур UserUrlRes с парами "короткий-оригинальный" URL.
	//   error - ошибка при получении данных из хранилища.
	UserUrls(ctx context.Context, userID int) (models.ApiUserUrlsRes, error)

	// Delete удаляет одну или несколько коротких ссылок по их идентификаторам.
	//
	// Параметры:
	//   ctx - контекст запроса.
	//   req - список ID коротких ссылок для удаления (ApiDeleteReq).
	//   userID - идентификатор пользователя, выполняющего удаление.
	//
	// Возвращает:
	//   []string - список удалённых ссылок.
	//   error - ошибка, если удаление не удалось выполнить.
	Delete(ctx context.Context, req models.ApiDeleteReq, userID int) ([]string, error)

	// InternalStats — получение внутренней статистики сервиса.
	// Возвращает агрегированные данные о работе системы, такие как количество сокращённых ссылок,
	// активных пользователей или другие метрики, специфичные для реализации.
	//
	// Параметры:
	//   ctx - контекст запроса.
	//
	// Возвращает:
	//   models.ApiInternalStatsRes - структура с данными внутренней статистики.
	//   error - ошибка при выполнении запроса к хранилищу данных.
	InternalStats(ctx context.Context) (models.ApiInternalStatsRes, error)

	// Ping проверяет работоспособность сервиса и его способность взаимодействовать с зависимостями (БД).
	//
	// Параметры:
	//   ctx - контекст запроса.
	Ping(ctx context.Context) error
}
