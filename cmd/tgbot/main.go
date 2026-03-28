package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/app"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/config"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/infra/broker"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/infra/httpclient"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/infra/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

const (
	TradingSignalsQueue = "trading-signals-queue"
	PnlReportsQueue     = "pnl-reports-queue"
	SystemQueue         = "system-queue"
)

func main() {
	os.Exit(run())
}

func run() int {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := logging.NewZapLogger("prod")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		return 1
	}
	defer func() {
		_ = logger.Sync() //nolint:errcheck // zap sync may fail on stdout
	}()

	goEnv := os.Getenv("GO_ENV")
	if goEnv != "prod" {
		if err := godotenv.Load(); err != nil {
			logger.Error("failed to load env", "error", err)
			return 1
		}
	}

	cfg, err := config.LoadFromEnv()
	if err != nil {
		logger.Error("failed to init config", "error", err)
		return 1
	}

	botAPI, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		logger.Error("failed to init crypto-knight telegram bot", "error", err)
		return 1
	}
	botAPI.Debug = false
	logger.Info("crypto-knight telegram bot started", "username", botAPI.Self.UserName)

	client := httpclient.New(cfg.APIBaseURL, time.Duration(cfg.HTTPTimeoutSeconds)*time.Second)

	brokerConn, err := broker.NewConnection(cfg.RmqURL)
	if err != nil {
		logger.Error("failed to init broker", "error", err)
		return 1
	}
	defer brokerConn.Close()
	logger.Info("broker started")

	err = brokerConn.DeclareQueue(TradingSignalsQueue)
	if err != nil {
		logger.Error("initialization failed", "type", TradingSignalsQueue, "error", err)
		return 1
	}
	logger.Info("initialized", "type", TradingSignalsQueue)

	err = brokerConn.DeclareQueue(PnlReportsQueue)
	if err != nil {
		logger.Error("initialization failed", "type", PnlReportsQueue, "error", err)
		return 1
	}
	logger.Info("initialized", "type", PnlReportsQueue)

	err = brokerConn.DeclareQueue(SystemQueue)
	if err != nil {
		logger.Error("initialization failed", "type", SystemQueue, "error", err)
		return 1
	}
	logger.Info("initialized", "type", SystemQueue)

	appl := app.NewApp(botAPI, cfg, client, brokerConn)

	done := make(chan error, 1)
	go func() {
		done <- appl.Run(ctx)
	}()

	select {
	case <-quit:
		logger.Info("graceful shutdown")
		cancel()
		if err := <-done; err != nil && !errors.Is(err, context.Canceled) {
			logger.Error("app exit error", "error", err)
		}
		return 0
	case err := <-done:
		if err != nil {
			logger.Error("app failed", "error", err)
			return 1
		}
		return 0
	}
}
