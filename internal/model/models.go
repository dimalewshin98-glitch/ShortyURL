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
