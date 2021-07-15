package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Harzu/rebrain-eventbus-webinar/kafka/internal/constants"
)

func New(logLevel string) (*zap.Logger, error) {
	zapLevel, err := level(logLevel).ToZap()
	if err != nil {
		return nil, err
	}

	zapConfig := zap.NewProductionConfig()
	zapConfig.Level.SetLevel(zapLevel)

	zapConfig.DisableStacktrace = true
	zapConfig.EncoderConfig.TimeKey = "timestamp"
	zapConfig.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	if _, err := zap.RedirectStdLogAt(logger, zapcore.DebugLevel); err != nil {
		return nil, err
	}

	return logger.With(zap.String("app", constants.ServiceName)), nil
}
