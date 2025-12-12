package command

import "github.com/google/uuid"

type GetURLByHashCommand struct {
	Hash string
}

type CreateURLEntryCommand struct {
	CorrelationID *string
	OriginalURL   string
	UserID        *uuid.UUID
}

type CreateBatchURLEntryCommand struct {
	Entries []CreateURLEntryCommand
}

type GetUrlsByUserIDCommand struct {
	UserID uuid.UUID
}
