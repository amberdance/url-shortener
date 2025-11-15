package errs

type UnauthorizedError string

func (e UnauthorizedError) Error() string {
	return string(e)
}

func (UnauthorizedError) Id() string {
	return "unauthorized"
}
