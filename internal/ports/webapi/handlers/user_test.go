package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	usecase "github.com/amberdance/url-shortener/internal/app/usecase/url"
	"github.com/amberdance/url-shortener/internal/domain/contracts"
	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/amberdance/url-shortener/internal/domain/repository"
	"github.com/amberdance/url-shortener/internal/infrastructure/helpers"
	infr "github.com/amberdance/url-shortener/internal/infrastructure/repository/url"
	"github.com/amberdance/url-shortener/internal/infrastructure/storage"
	"github.com/amberdance/url-shortener/internal/ports/webapi/dto"
	ports "github.com/amberdance/url-shortener/internal/ports/webapi/helpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	userUrlsEndpoint = "/urls"
)

type handlerWrapper struct {
	repository repository.URLRepository
	handler    *UserHandler
	host       string
}

func buildUserHandler() *handlerWrapper {
	repo := infr.NewInMemoryURLRepository(storage.NewInMemoryStorage())
	host := "http://localhost:8080/"
	return &handlerWrapper{
		repository: repo,
		host:       host,
		handler:    NewUserHandler(host, usecase.NewGetURLsByUserIDUseCase(repo)),
	}
}

func Test_When_UserHasUrls_Then_URLsReturned(t *testing.T) {
	h := buildUserHandler()
	id, ctx := generateUuidWithContext(t.Context())
	urls := seedUrls(h.repository, &id)
	req := httptest.NewRequest(http.MethodGet, userUrlsEndpoint, nil).WithContext(ctx)
	w := httptest.NewRecorder()

	req.Header.Set("user_id", "value")

	h.handler.Routes().ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()

	response, _ := io.ReadAll(res.Body)

	var dtos []dto.UserURLsResponse
	err := json.Unmarshal(response, &dtos)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, len(urls), len(dtos))

	sort.Slice(dtos, func(i, j int) bool {
		return dtos[i].OriginalURL < dtos[j].OriginalURL
	})
	sort.Slice(urls, func(i, j int) bool {
		return urls[i].OriginalURL < urls[j].OriginalURL
	})

	for i := range urls {
		assert.Equal(t, urls[i].OriginalURL, dtos[i].OriginalURL)
		assert.Equal(t, ports.FormatFullURL(h.host, urls[i].Hash), dtos[i].ShortURL)
	}
}

func Test_When_UserDoesNotHasUrls_Then_204HttpCodeReturned(t *testing.T) {
	h := buildUserHandler()
	_, ctx := generateUuidWithContext(t.Context())
	req := httptest.NewRequest(http.MethodGet, userUrlsEndpoint, nil).WithContext(ctx)
	w := httptest.NewRecorder()

	req.Header.Set("user_id", "YmFhYzY0NzEtZGFhZS00NGY3LWE0M2QtODhmYmY0YTU3Mzlm.V9gkZuW7x7qAf8aG3BoBYcVWwKd6KWClgcxnQimAlnA")

	h.handler.Routes().ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusNoContent, res.StatusCode)
}

func Test_When_HeaderNotPresent_Then_401HttpCodeReturned(t *testing.T) {
	h := buildUserHandler()
	req := httptest.NewRequest(http.MethodGet, userUrlsEndpoint, nil)
	w := httptest.NewRecorder()

	h.handler.Routes().ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func Test_When_InvalidSignatureProvided_Then_401HttpCodeReturned(t *testing.T) {
	h := buildUserHandler()

	badCookie := &http.Cookie{
		Name:  "user_id",
		Value: "broken.token.value",
		Path:  "/",
	}

	req := httptest.NewRequest(http.MethodGet, userUrlsEndpoint, nil)
	req.AddCookie(badCookie)

	w := httptest.NewRecorder()
	h.handler.Routes().ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func generateUuidWithContext(ctx context.Context) (uuid.UUID, context.Context) {
	id := uuid.New()
	c := context.WithValue(ctx, contracts.UserCtxKey, id.String())

	return id, c
}

func seedUrls(r repository.URLRepository, userId *uuid.UUID) []*model.URL {
	urls := make([]*model.URL, 0, 10)

	for i := 0; i < 10; i++ {
		m, _ := model.NewURL(fmt.Sprintf("https://original-%d.ru", i), helpers.GenerateHash(), nil, nil)
		m.UserID = userId
		urls = append(urls, m)
		_ = r.Create(context.TODO(), urls[i])
	}

	return urls
}
