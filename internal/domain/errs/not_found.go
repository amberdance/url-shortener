package errs

type NotFoundError string

func (e NotFoundError) Error() string { return string(e) }

func (NotFoundError) ID() string { return "not_found" }
