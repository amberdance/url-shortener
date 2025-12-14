package command

import "github.com/google/uuid"

type GetURLByHashCommand struct {
	Hash string
}

type CreateURLEntryCommand struct {
	CorrelationID *string
	OriginalURL   string
	UserID        *uuid.UUID
	RequestID     uuid.UUID
}

type CreateBatchURLEntryCommand struct {
	Commands  []CreateURLEntryCommand
	RequestID uuid.UUID
}

type GetUrlsByUserIDCommand struct {
	UserID uuid.UUID
}
