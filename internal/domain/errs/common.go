package errs

type ValidationError string

func (e ValidationError) Error() string {
	return string(e)
}

func (ValidationError) Id() string { return "validation_error" }
