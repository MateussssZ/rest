package appLogger

import (
	"context"
	"log/slog"
)

type IAppLogger interface {
	Debug(ctx context.Context, msg string, attrs ...any)
	Info(ctx context.Context, msg string, attrs ...any)
	Error(ctx context.Context, err error, attrs ...any)
}

type AppLogger struct {
	logger *stdoutWriter
}

func NewAppLogger(level slog.Level) *AppLogger {

	return &AppLogger{
		logger: newStdoutWriter(level),
	}
}

func (l *AppLogger) Debug(ctx context.Context, msg string, attrs ...any) {
	l.logger.debugContext(ctx, msg, attrs...)
}

func (l *AppLogger) Info(ctx context.Context, msg string, attrs ...any) {
	l.logger.infoContext(ctx, msg, attrs...)
}

func (l *AppLogger) Error(ctx context.Context, err error, attrs ...any) {
	l.logger.errorContext(ctx, err.Error(), attrs...)
}
