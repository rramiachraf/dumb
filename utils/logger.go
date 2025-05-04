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

func (l *Logger) Errorf(f string, args ...any) {
	l.slog.Error(fmt.Sprintf(f, args...))
}

func (l *Logger) Error(e string) {
	l.Errorf("%s", e)
}

func (l *Logger) Infof(f string, args ...any) {
	l.slog.Info(fmt.Sprintf(f, args...))
}

func (l *Logger) Info(m string) {
	l.Infof("%s", m)
}

func (l *Logger) Fatalf(f string, args ...any) {
	l.Errorf(f, args...)
	os.Exit(1)
}

func (l *Logger) Fatal(e string) {
	l.Fatalf("%s", e)
}
