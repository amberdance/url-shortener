package errs

type UnauthorizedError string

func (e UnauthorizedError) Error() string {
	return string(e)
}

func (UnauthorizedError) ID() string {
	return "unauthorized"
}
