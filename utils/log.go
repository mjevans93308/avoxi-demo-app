package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger returns a zap sugared logger instance with ISO8601 timestamps
// Using the production config level gives us JSON encoding for our log messages
// as well as stacktraces for ERROR-level logs and caller id
func InitLogger() *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := config.Build()
	return logger.Sugar()
}
