package models

// ApiShortenReq представляет собой структуру запроса на создание одной короткой ссылки через API.
type ApiShortenReq struct {
	// URL — исходный длинный URL-адрес, который необходимо сократить.
	URL string `json:"url"`
}

// ApiShortenRes представляет собой структуру ответа при успешном создании одной короткой ссылки.
type ApiShortenRes struct {
	// Result — созданная короткая ссылка.
	Result string `json:"result"`
}

// BatchItemReq описывает один элемент в пакете запросов на сокращение URL.
type BatchItemReq struct {
	// CorrelationID — уникальный идентификатор элемента для сопоставления запроса и ответа.
	CorrelationID string `json:"correlation_id"`
	// OriginalURL — исходный длинный URL-адрес.
	OriginalURL string `json:"original_url"`
}

// ApiShortenBatchReq представляет собой срез элементов BatchItemReq для пакетной обработки.
type ApiShortenBatchReq []BatchItemReq

// BatchItemRes описывает результат сокращения для одного элемента пакета.
type BatchItemRes struct {
	// CorrelationID — идентификатор, пришедший в запросе.
	CorrelationID string `json:"correlation_id"`
	// ShortURL — результат сокращения.
	ShortURL string `json:"short_url"`
}

// ApiShortenBatchRes представляет собой срез результатов BatchItemRes.
type ApiShortenBatchRes []BatchItemRes

// UserUrlRes содержит информацию об одной паре "короткий - оригинальный" URL пользователя.
type UserUrlRes struct {
	// ShortURL — сокращенная версия ссылки.
	ShortURL string `json:"short_url"`
	// OriginalURL — оригинальная длинная ссылка.
	OriginalURL string `json:"original_url"`
}

// ApiUserUrlsRes представляет собой срез информации о всех ссылках пользователя.
type ApiUserUrlsRes []UserUrlRes

// ApiDeleteReq представляет собой список ID коротких ссылок для удаления.
type ApiDeleteReq []string

// RepoDeleteMessage используется для передачи команды на удаление в слой хранения данных.
type RepoDeleteMessage struct {
	// ShortURL — ID короткой ссылки, которую нужно удалить.
	ShortURL string
	// UserID — ID пользователя-владельца ссылки.
	UserID int
}

// ApiDeleteReq представляет собой список удаленных url
type ApiDeleteRes []string

// AuditData содержит данные для аудита действий пользователя.
type AuditData struct {
	// TS — временная метка (timestamp) совершения действия.
	TS int64 `json:"ts"`
	// Action — тип совершенного действия.
	Action string `json:"action"`
	// UserID — идентификатор пользователя.
	UserID string `json:"user_id"`
	// URL — URL-адрес, над которым было совершено действие.
	URL string `json:"url"`
}
