package dto

type ShortURLRequest struct {
	URL string `json:"url" validate:"required"`
}
type ShortURLResponse struct {
	Result string `json:"result"`
}
