package models

type ApiShortenReq struct {
	URL string `json:"url"`
}

type ApiShortenRes struct {
	Result string `json:"result"`
}

type ApiShortenBatchReq []BatchItemReq

type BatchItemReq struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ApiShortenBatchRes []BatchItemRes

type BatchItemRes struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type ApiUserUrlsRes []UserUrlRes

type UserUrlRes struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type ApiDeleteReq []string

type ApiDeleteRes []string

type RepoDeleteMessage struct {
	ShortURL string
	UserID   int
}

type AuditData struct {
	TS     int64  `json:"ts"`
	Action string `json:"action"`
	UserID string `json:"user_id"`
	URL    string `json:"url"`
}
