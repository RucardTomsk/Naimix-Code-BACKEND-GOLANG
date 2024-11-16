package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// NewLogger returns preconfigured *zap.Logger.
func NewLogger() *zap.Logger {
	logger, _ := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "timestamp",
			EncodeTime: zapcore.EpochTimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}.Build()

	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger) // flush output if any

	return logger.With(
		zap.String("env", os.Getenv("ENV")),
		zap.String("system", os.Getenv("SYSTEM")),
		zap.String("inst", os.Getenv("INST")),
	)
}
