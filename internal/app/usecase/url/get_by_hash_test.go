package url

import (
	"context"
	"testing"

	"github.com/amberdance/url-shortener/internal/app/command"
	"github.com/amberdance/url-shortener/internal/infrastructure/repository/url"
	"github.com/amberdance/url-shortener/internal/infrastructure/storage"
	"github.com/stretchr/testify/suite"
)

type GetByHashUseCaseTestSuite struct {
	suite.Suite
	ctx context.Context

	useCase       GetByHashUseCase
	memoryStorage storage.InMemoryStorage
}

func (s *GetByHashUseCaseTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.memoryStorage = *storage.NewInMemoryStorage()
	s.useCase = NewGetByHashUseCase(url.NewInMemoryURLRepository(&s.memoryStorage))
}

func (s *GetByHashUseCaseTestSuite) TearDownTest() {
	clear(s.memoryStorage.Data)
}

func TestGetByHashUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(GetByHashUseCaseTestSuite))
}

func (s *GetByHashUseCaseTestSuite) TestEntryNotFound() {
	hash := "some-hash"
	entry, err := s.useCase.Run(s.ctx, command.GetURLByHashCommand{
		Hash: hash,
	})
	s.Error(err)
	s.Empty(entry)
}

func (s *GetByHashUseCaseTestSuite) TestGetByHashSuccess() {
	expected := createTestURLEntry()
	s.memoryStorage.Data[expected.ID] = expected
	actual, err := s.useCase.Run(s.ctx, command.GetURLByHashCommand{
		Hash: expected.Hash,
	})
	s.NoError(err)
	s.Equal(expected, actual)
}
