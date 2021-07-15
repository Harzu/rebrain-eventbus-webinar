package logger

import (
	"errors"

	"go.uber.org/zap/zapcore"
)

type level string

func (l level) ToZap() (zapcore.Level, error) {
	switch l {
	case "DEBUG":
		return zapcore.DebugLevel, nil
	case "INFO":
		return zapcore.InfoLevel, nil
	case "WARN", "WARNING":
		return zapcore.WarnLevel, nil
	case "ERROR":
		return zapcore.ErrorLevel, nil
	case "DPANIC":
		return zapcore.DPanicLevel, nil
	case "PANIC":
		return zapcore.PanicLevel, nil
	case "FATAL":
		return zapcore.FatalLevel, nil
	}

	return 0, errors.New("invalid log level name")
}
