// Package logging provides a Zap-backed implementation of ports.Logger.
package logging

import (
	"fmt"

	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/ports"
	"go.uber.org/zap"
)

// ZapLogger implements ports.Logger using zap's sugared logger.
type ZapLogger struct {
	sugar *zap.SugaredLogger
}

// NewZapLogger builds a logger using production or development zap config from env.
func NewZapLogger(env string) (ports.Logger, error) {
	var cfg zap.Config
	if env == "prod" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("zap build: %w", err)
	}

	return &ZapLogger{sugar: logger.Sugar()}, nil
}

// Sync flushes any buffered log entries.
func (l *ZapLogger) Sync() error {
	if err := l.sugar.Sync(); err != nil {
		return fmt.Errorf("zap sync: %w", err)
	}
	return nil
}

// Info logs at info level with optional structured fields.
func (l *ZapLogger) Info(msg string, keysAndValues ...any) {
	l.sugar.Infow(msg, keysAndValues...)
}

// Error logs at error level with optional structured fields.
func (l *ZapLogger) Error(msg string, keysAndValues ...any) {
	l.sugar.Errorw(msg, keysAndValues...)
}

// Debug logs at debug level with optional structured fields.
func (l *ZapLogger) Debug(msg string, keysAndValues ...any) {
	l.sugar.Debugw(msg, keysAndValues...)
}
