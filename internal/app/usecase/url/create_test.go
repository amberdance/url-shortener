package url

import (
	"context"
	"testing"
	"time"

	"github.com/amberdance/url-shortener/internal/app/command"
	"github.com/amberdance/url-shortener/internal/domain/errs"
	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/amberdance/url-shortener/internal/infrastructure/repository/url"
	"github.com/amberdance/url-shortener/internal/infrastructure/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type CreateUseCaseTestSuite struct {
	suite.Suite
	ctx context.Context

	useCase       CreateUseCase
	memoryStorage storage.InMemoryStorage
}

func (s *CreateUseCaseTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.memoryStorage = *storage.NewInMemoryStorage()
	s.useCase = NewCreateUseCase(url.NewInMemoryURLRepository(&s.memoryStorage))
}

func (s *CreateUseCaseTestSuite) TearDownTest() {
	clear(s.memoryStorage.Data)
}

func TestCreateURLUseCaseSuite(t *testing.T) {
	suite.Run(t, new(CreateUseCaseTestSuite))
}

func (s *CreateUseCaseTestSuite) TestSuccess() {
	entry, err := s.useCase.Run(s.ctx, command.CreateURLEntryCommand{
		OriginalURL: "https://test.com",
	})
	s.NoError(err)

	expected := s.memoryStorage.Data[entry.ID]
	s.Equal(expected, entry)
}

func (s *CreateUseCaseTestSuite) TestDuplicateEntry() {
	entry := createTestURLEntry()
	entry.OriginalURL = "duplicate"
	s.memoryStorage.Data[entry.ID] = entry

	_, err := s.useCase.Run(s.ctx, command.CreateURLEntryCommand{OriginalURL: "duplicate"})
	var dep errs.DuplicateEntryError
	s.ErrorAs(err, &dep)
}

func createTestURLEntry() *model.URLEntry {
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
