package errs

var (
	ErrNotFound      = NotFoundError("не найден запрашиваемый ресурс")
	ErrDuplicate     = DuplicateEntryError("запись уже существует")
	ErrUnauthorized  = UnauthorizedError("unauthorized")
	ErrInvalidUserID = UnauthorizedError("invalid user id")
	ErrInvalidURI    = ValidationError("некорректный URI")
	ErrIncorrectURL  = ValidationError("не удалось сформировать ссылку")
	ErrEmptyURLSet   = ValidationError("не передано ни одного url")
)
