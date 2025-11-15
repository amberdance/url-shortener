package errs

type InvalidArgumentError string

func (e InvalidArgumentError) Error() string {
	return string(e)
}

func (InvalidArgumentError) ID() string { return "invalid_argument" }
