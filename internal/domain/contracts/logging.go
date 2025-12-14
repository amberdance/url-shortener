package contracts

//go:generate mockgen -source=$GOFILE -destination=../../mocks/logger_mock.go -package=mocks
type Logger interface {
	Debug(message string, args ...any)
	Info(message string, args ...any)
	Error(message string, args ...any)
	Close() error
}
