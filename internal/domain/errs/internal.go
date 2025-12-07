package errs

type InternalError string

func (e InternalError) Error() string {
	return string(e)
}

func (InternalError) ID() string { return "internal_error" }
