package service

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/amberdance/url-shortener/internal/app"
	"github.com/amberdance/url-shortener/internal/ports/webapi/helpers"
)

type URLShortenerService struct {
	storage app.Storage
}

func NewURLShortenerService(storage app.Storage) *URLShortenerService {
	return &URLShortenerService{storage: storage}
}

func (s *URLShortenerService) CreateShortURL(ctx context.Context, original string) (string, error) {
	original = strings.TrimSpace(original)
	if original == "" {
		return "", errors.New("empty URL")
	}

	if _, err := url.ParseRequestURI(original); err != nil {
		return "", errors.New("invalid URL format")
	}

	id := helpers.GenerateShortID()
	if err := s.storage.Save(ctx, id, original); err != nil {
		return "", err
	}

	return id, nil
}

func (s *URLShortenerService) ResolveURL(ctx context.Context, id string) (string, error) {
	entry, ok := s.storage.Get(ctx, id)
	if !ok {
		return "", errors.New("entry not found")
	}
	return entry, nil
}
