package log

import (
	"context"
	"go.uber.org/zap"
)

type ctxLogger struct{}

// ContextWithLogger adds logger to context.
func ContextWithLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, l)
}

// LoggerFromContext returns logger from context, if exists.
// Creates new one otherwise.
func LoggerFromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxLogger{}).(*zap.Logger); ok {
		return l
	}

	return NewLogger()
}
