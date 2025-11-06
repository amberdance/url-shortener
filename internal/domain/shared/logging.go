package shared

type Logger interface {
	Debug(message string, args ...any)
	Info(message string, args ...any)
	Error(message string, args ...any)
	Close() error
}
