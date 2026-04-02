package models

type ApiRequest struct {
	URL string `json:"url"`
}

type ApiResponse struct {
	Result string `json:"result"`
}
