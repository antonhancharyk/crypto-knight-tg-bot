package main

import (
	"context"
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

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	logger, err := logging.NewZapLogger("prod")
	if err != nil {
		panic(`failed to init logger: ` + err.Error())
	}
	defer func() {
		_ = logger.Sync()
	}()

	GO_ENV := os.Getenv("GO_ENV")
	if GO_ENV != "prod" {
		err := godotenv.Load()
		if err != nil {
			logger.Error("failed to load env", "error", err)
			os.Exit(1)
		}
	}

	cfg, err := config.LoadFromEnv()
	if err != nil {
		logger.Error("failed to init config", "error", err)
		os.Exit(1)
	}

	botAPI, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		logger.Error("failed to init crypto-knight telegram bot", "error", err)
		os.Exit(1)
	}
	botAPI.Debug = false
	logger.Info("crypto-knight telegram bot started", "username", botAPI.Self.UserName)

	client := httpclient.New(cfg.APIBaseURL, time.Duration(cfg.HTTPTimeoutSeconds)*time.Second)

	brokerConn, err := broker.NewConnection(cfg.RmqURL)
	if err != nil {
		logger.Error("failed to init broker", "error", err)
		os.Exit(1)
	}
	defer brokerConn.Close()
	logger.Info("broker started")

	err = brokerConn.DeclareQueue("tg")
	if err != nil {
		logger.Error("initialization failed", "type", "queue_tg", "error", err)
		os.Exit(1)
	}
	logger.Info("initialized", "type", "queue_tg")

	appl := app.NewApp(botAPI, cfg, client, brokerConn)
	if err := appl.Run(ctx); err != nil {
		logger.Error("failed to run app", "error", err)
		os.Exit(1)
	}

	<-quit
	logger.Info("graceful shutdown")
	cancel()
}
