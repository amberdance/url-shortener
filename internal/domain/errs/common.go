package errs

type ValidationError string

func (e ValidationError) Error() string {
	return string(e)
}

func (ValidationError) ID() string { return "validation_error" }
