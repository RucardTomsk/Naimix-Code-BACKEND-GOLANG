package log

import (
	"fmt"
	"go.uber.org/zap"
)
import (
	gormlogger "gorm.io/gorm/logger"
)

// ZapWriter is a custom logger that implements the gormlogger.Writer interface.
type ZapWriter struct {
	logger *zap.Logger
}

// NewZapWriter creates a new ZapWriter.
func NewZapWriter(logger *zap.Logger) gormlogger.Writer {
	return &ZapWriter{logger: logger}
}

// Printf implements the Printf method of the Writer interface.
func (zw *ZapWriter) Printf(format string, args ...interface{}) {
	// Use zap's Sprintf function to format the log message
	msg := fmt.Sprintf(format, args...)

	// Log the formatted message using Zap's Info method
	zw.logger.Info(msg)
}
