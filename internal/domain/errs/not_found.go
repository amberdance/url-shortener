package errs

type NotFoundError string

func (e NotFoundError) Error() string {
	return string(e)
}

func (NotFoundError) Id() string { return "not_found" }
