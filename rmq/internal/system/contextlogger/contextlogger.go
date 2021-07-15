package contextlogger

import (
	"context"

	"go.uber.org/zap"
)

const loggerKey = "LOGGER"

func Enrich(ctx context.Context, l *zap.Logger) context.Context {
	//nolint
	return context.WithValue(ctx, loggerKey, l)
}

func Fetch(ctx context.Context) *zap.Logger {
	l, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		return zap.NewNop()
	}

	return l
}
