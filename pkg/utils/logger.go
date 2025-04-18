package utils

import "go.uber.org/zap"

func NewLogger(level string) *zap.Logger {
	cfg := zap.NewProductionConfig()
	if level == "debug" {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	logger, _ := cfg.Build()
	return logger
}
