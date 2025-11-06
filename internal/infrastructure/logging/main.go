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
	file       *os.File
	termLogger *slog.Logger
	fileLogger *slog.Logger
}

var _ shared.Logger = (*Logger)(nil)

func NewLogger() *Logger {
	cfg := config.NewConfig()
	f, err := os.OpenFile("./logs/app.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln(err)
	}

	level := slog.LevelDebug

	switch cfg.LogLevel {
	case "info":
		level = slog.LevelInfo
	case "error":
		level = slog.LevelError
	}

	fileHandler := tint.NewHandler(f, &tint.Options{
		Level:      level,
		TimeFormat: time.DateTime,
		NoColor:    true,
	})

	stdoutHandler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      level,
		TimeFormat: time.DateTime,
		NoColor:    false,
	})

	return &Logger{
		file:       f,
		fileLogger: slog.New(fileHandler),
		termLogger: slog.New(stdoutHandler),
	}
}

func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

func (l *Logger) Info(msg string, args ...any) {
	if len(args) > 0 {
		l.fileLogger.Info(msg, args...)
		l.termLogger.Info(msg, args...)
	} else {
		l.fileLogger.Info(msg)
		l.termLogger.Info(msg)
	}
}

func (l *Logger) Debug(msg string, args ...any) {
	if len(args) > 0 {
		l.fileLogger.Debug(msg, args...)
		l.termLogger.Debug(msg, args...)
	} else {
		l.fileLogger.Debug(msg)
		l.termLogger.Debug(msg)
	}
}

func (l *Logger) Error(msg string, args ...any) {
	if len(args) > 0 {
		l.fileLogger.Error(msg, args...)
		l.termLogger.Error(msg, args...)
	} else {
		l.fileLogger.Error(msg)
		l.termLogger.Error(msg)
	}
}
