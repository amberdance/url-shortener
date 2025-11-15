package errs

type InvalidArgumentError string

func (e InvalidArgumentError) Error() string {
	return string(e)
}

func (InvalidArgumentError) Id() string { return "invalid_argument" }
