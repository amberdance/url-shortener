package dto

type ShortURLRequest struct {
	CorrelationID *string `json:"correlation_id"`
	URL           string  `json:"url" validate:"required"`
}
type ShortURLResponse struct {
	URL string `json:"result"`
}

type BatchShortenURLRequest struct {
	URL           string `json:"original_url" validate:"required"`
	CorrelationID string `json:"correlation_id" validate:"required"`
}

type BatchShortenURLResponse struct {
	URL           string `json:"short_url"`
	CorrelationID string `json:"correlation_id"`
}

type UserURLsResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
