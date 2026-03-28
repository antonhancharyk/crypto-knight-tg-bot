package logging

import (
	"fmt"

	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/ports"
	"go.uber.org/zap"
)

type ZapLogger struct {
	sugar *zap.SugaredLogger
}

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

func (l *ZapLogger) Sync() error {
	if err := l.sugar.Sync(); err != nil {
		return fmt.Errorf("zap sync: %w", err)
	}
	return nil
}

func (l *ZapLogger) Info(msg string, keysAndValues ...any) {
	l.sugar.Infow(msg, keysAndValues...)
}

func (l *ZapLogger) Error(msg string, keysAndValues ...any) {
	l.sugar.Errorw(msg, keysAndValues...)
}

func (l *ZapLogger) Debug(msg string, keysAndValues ...any) {
	l.sugar.Debugw(msg, keysAndValues...)
}
