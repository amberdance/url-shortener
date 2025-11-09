package url_test

import (
	"context"
	"testing"

	"github.com/amberdance/url-shortener/internal/app/command"
	urlusecase "github.com/amberdance/url-shortener/internal/app/usecase/url"
	"github.com/amberdance/url-shortener/internal/infrastructure/repository/url"
	"github.com/stretchr/testify/assert"
)

func TestGetByHashUseCase_Run_Success(t *testing.T) {
	repo := url.NewInMemoryRepository()
	create := urlusecase.NewCreateUrlUseCase(repo)
	get := urlusecase.NewGetByHashUseCase(repo)
	cmd := command.CreateURLEntryCommand{OriginalURL: "https://hard2code.ru"}

	m, err := create.Run(context.Background(), cmd)
	assert.NoError(t, err)

	found, err := get.Run(context.Background(), command.GetURLByHashCommand{Hash: m.Hash})
	assert.NoError(t, err)
	assert.Equal(t, m, found)
}

func TestGetByHashUseCase_Run_NotFound(t *testing.T) {
	repo := url.NewInMemoryRepository()
	get := urlusecase.NewGetByHashUseCase(repo)

	_, err := get.Run(context.Background(), command.GetURLByHashCommand{Hash: "none"})
	assert.Error(t, err)
}
