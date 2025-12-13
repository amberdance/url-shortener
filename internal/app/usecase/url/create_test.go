package url_test

import (
	"context"
	"testing"
	"time"

	"github.com/amberdance/url-shortener/internal/app/command"
	uc "github.com/amberdance/url-shortener/internal/app/usecase/url"
	"github.com/amberdance/url-shortener/internal/domain/errs"
	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/amberdance/url-shortener/internal/infrastructure/repository/url"
	"github.com/amberdance/url-shortener/internal/infrastructure/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type CreateURLUseCaseTestSuite struct {
	suite.Suite
	ctx context.Context

	useCase       uc.CreateUseCase
	memoryStorage storage.InMemoryStorage
}

func (s *CreateURLUseCaseTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.memoryStorage = *storage.NewInMemoryStorage()
	s.useCase = uc.NewCreateURLUseCase(url.NewInMemoryURLRepository(&s.memoryStorage))
}

func (s *CreateURLUseCaseTestSuite) TearDownTest() {
	clear(s.memoryStorage.Data)
}

func TestCreateURLUseCaseSuite(t *testing.T) {
	suite.Run(t, new(CreateURLUseCaseTestSuite))
}

func (s *CreateURLUseCaseTestSuite) TestSuccess() {
	entry, err := s.useCase.Run(s.ctx, command.CreateURLEntryCommand{
		OriginalURL: "https://test.com",
	})
	s.NoError(err)

	expected := s.memoryStorage.Data[entry.ID]
	s.Equal(expected, entry)
}

func (s *CreateURLUseCaseTestSuite) TestDuplicateEntry() {
	entry := s.createURLEntry()
	entry.OriginalURL = "duplicate"
	s.memoryStorage.Data[entry.ID] = entry

	_, err := s.useCase.Run(s.ctx, command.CreateURLEntryCommand{OriginalURL: "duplicate"})
	var dep errs.DuplicateEntryError
	s.ErrorAs(err, &dep)
}

func (s *CreateURLUseCaseTestSuite) createURLEntry() *model.URLEntry {
	userID := uuid.New()
	entryID := uuid.New()

	return &model.URLEntry{
		ID:          entryID,
		OriginalURL: "https://example.com",
		Hash:        "hash",
		CreatedAt:   time.Now(),
		UserID:      &userID,
	}
}
