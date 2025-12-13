package url

import (
	"context"
	"fmt"
	"testing"

	"github.com/amberdance/url-shortener/internal/app/command"
	"github.com/amberdance/url-shortener/internal/domain/model"
	"github.com/amberdance/url-shortener/internal/infrastructure/repository/url"
	"github.com/amberdance/url-shortener/internal/infrastructure/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type BatchCreateURLUseCaseTestSuite struct {
	suite.Suite
	ctx context.Context

	useCase       BatchCreateURLUseCase
	memoryStorage storage.InMemoryStorage
}

func (s *BatchCreateURLUseCaseTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.memoryStorage = *storage.NewInMemoryStorage()
	s.useCase = NewBatchCreateURLUseCase(url.NewInMemoryURLRepository(&s.memoryStorage))
}

func (s *BatchCreateURLUseCaseTestSuite) TearDownTest() {
	clear(s.memoryStorage.Data)
}

func TestBatchCreateURLUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(BatchCreateURLUseCaseTestSuite))
}

func (s *BatchCreateURLUseCaseTestSuite) TestCreateBatchSuccess() {
	commands := make([]command.CreateURLEntryCommand, 3)
	someID := uuid.New()
	requestID := someID.String()
	for i := range commands {
		commands[i] = command.CreateURLEntryCommand{
			CorrelationID: &requestID,
			OriginalURL:   fmt.Sprintf("http://test-domain%d.com", i),
			UserID:        &someID,
		}
	}

	result, err := s.useCase.Run(s.ctx, command.CreateBatchURLEntryCommand{Commands: commands})
	s.NoError(err)

	for i, u := range result {
		s.Equal(u.OriginalURL, commands[i].OriginalURL)
	}
}

func (s *BatchCreateURLUseCaseTestSuite) TestCreateBatchDuplicate() {
	urls := make([]*model.URLEntry, 3)
	for i := range urls {
		urls[i] = createTestURLEntry()
		s.memoryStorage.Data[urls[i].ID] = urls[i]
	}

	commands := make([]command.CreateURLEntryCommand, 3)
	for i := 0; i < 2; i++ {
		commands[i] = command.CreateURLEntryCommand{OriginalURL: urls[i].OriginalURL}
	}

	commands[len(commands)-1] = command.CreateURLEntryCommand{OriginalURL: "some-url"}
	result, err := s.useCase.Run(s.ctx, command.CreateBatchURLEntryCommand{Commands: commands})
	s.Error(err)
	s.Nil(result)
}
