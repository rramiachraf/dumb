package utils

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

type Logger struct {
	slog *slog.Logger
}

func NewLogger(w io.WriteCloser) *Logger {
	handler := slog.NewTextHandler(w, &slog.HandlerOptions{})
	sl := slog.New(handler)

	return &Logger{slog: sl}
}

func (l *Logger) Error(f string, args ...any) {
	l.slog.Error(fmt.Sprintf(f, args...))
}

func (l *Logger) Info(f string, args ...any) {
	l.slog.Info(fmt.Sprintf(f, args...))
}

func (l *Logger) Fatal(f string, args ...any) {
	l.Error(f, args...)
	os.Exit(1)
}
