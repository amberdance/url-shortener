package errs

type DuplicateEntryError string

func (e DuplicateEntryError) Error() string {
	return string(e)
}

func (DuplicateEntryError) Id() string { return "duplicate_entry" }
