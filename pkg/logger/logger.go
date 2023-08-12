package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
)

func init() {
	cfg := zap.NewProductionConfig()
	cfg.DisableCaller = false
	cfg.DisableStacktrace = false
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	localLogger, err := cfg.Build()
	if err != nil {
		log.Fatal("logger init", err)
	}
	logger = localLogger
}

func Get() *zap.SugaredLogger {
	return logger.Sugar()
}
