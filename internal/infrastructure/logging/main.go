package logging

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/amberdance/url-shortener/internal/config"
	"github.com/amberdance/url-shortener/internal/domain/shared"
	"github.com/lmittmann/tint"
)

type Logger struct {
	file *os.File
	logs []*slog.Logger
}

var _ shared.Logger = (*Logger)(nil)

func NewLogger() *Logger {
	cfg := config.GetConfig()

	f, err := os.OpenFile("./logs/app.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	level := slog.LevelDebug
	switch cfg.LogLevel {
	case "info":
		level = slog.LevelInfo
	case "error":
		level = slog.LevelError
	}

	opts := &tint.Options{Level: level, TimeFormat: time.DateTime}

	return &Logger{
		file: f,
		logs: []*slog.Logger{
			slog.New(tint.NewHandler(os.Stdout, opts)),
			slog.New(tint.NewHandler(f, &tint.Options{
				Level:      level,
				TimeFormat: time.DateTime,
			})),
		},
	}
}

func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

func (l *Logger) log(fn func(*slog.Logger, string, ...any), msg string, args ...any) {
	for _, lg := range l.logs {
		fn(lg, msg, args...)
	}
}

func (l *Logger) Info(msg string, args ...any)  { l.log((*slog.Logger).Info, msg, args...) }
func (l *Logger) Debug(msg string, args ...any) { l.log((*slog.Logger).Debug, msg, args...) }
func (l *Logger) Error(msg string, args ...any) { l.log((*slog.Logger).Error, msg, args...) }
