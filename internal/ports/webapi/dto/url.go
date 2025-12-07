package dto

type ShortURLRequest struct {
	CorrelationID *string `json:"correlation_id" json:"correlation_id"`
	URL           string  `json:"url" validate:"required,url"`
}
type ShortURLResponse struct {
	URL string `json:"url"`
}

type BatchShortenURLRequest struct {
	CorrelationID string `json:"correlation_id" validate:"required"`
	OriginalURL   string `json:"original_url" validate:"required,url"`
}

type BatchShortenURLResponse struct {
	CorrelationID string `json:"correlation_id" validate:"required"`
	URL           string `json:"short_url"`
}
