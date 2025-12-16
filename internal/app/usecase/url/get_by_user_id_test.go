package url

import (
	"context"
	"testing"

	"github.com/amberdance/url-shortener/internal/app/command"
	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/amberdance/url-shortener/internal/infrastructure/repository/url"
	"github.com/amberdance/url-shortener/internal/infrastructure/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type GetURLsByUserIDUseCaseTestSuite struct {
	suite.Suite
	ctx context.Context

	useCase       GetURLsByUserIDUseCase
	memoryStorage storage.InMemoryStorage
}

func (s *GetURLsByUserIDUseCaseTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.memoryStorage = *storage.NewInMemoryStorage()
	s.useCase = NewGetURLsByUserIDUseCase(url.NewInMemoryURLRepository(&s.memoryStorage))
}

func (s *GetURLsByUserIDUseCaseTestSuite) TearDownTest() {
	clear(s.memoryStorage.Data)
}

func TestGetURLsByUserIDUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(GetURLsByUserIDUseCaseTestSuite))
}

func (s *GetURLsByUserIDUseCaseTestSuite) TestGetByUserIDSuccess() {
	urls := []*model.URLEntry{
		createTestURLEntry(),
		createTestURLEntry(),
		createTestURLEntry(),
	}
	userID := uuid.New()

	for _, u := range urls {
		u.UserID = &userID
		s.memoryStorage.Data[u.ID] = u
	}
	res, err := s.useCase.Run(s.ctx, command.GetUrlsByUserIDCommand{UserID: userID})
	s.NoError(err)
	s.Equal(urls, res)
}
