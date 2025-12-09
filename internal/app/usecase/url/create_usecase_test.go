package url_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amberdance/url-shortener/internal/app/command"
	urlusecase "github.com/amberdance/url-shortener/internal/app/usecase/url"
	"github.com/amberdance/url-shortener/internal/infrastructure/repository/url"
)

func TestCreateUseCase_Run_Success(t *testing.T) {
	uc := urlusecase.NewCreateURLUseCase(url.NewInMemoryURLRepository())
	cmd := command.CreateURLEntryCommand{
		OriginalURL: "https://hard2code.ru",
	}

	m, err := uc.Run(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, cmd.OriginalURL, m.OriginalURL)
	assert.NotEmpty(t, m.Hash)
}
