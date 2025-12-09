package command

type GetURLByHashCommand struct {
	Hash string
}

type CreateURLEntryCommand struct {
	CorrelationID *string
	OriginalURL   string
}

type CreateBatchURLEntryCommand struct {
	Entries []CreateURLEntryCommand
}
