package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/app"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/config"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/infra/httpclient"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/infra/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-quit
		cancel()
	}()

	err := godotenv.Load()

	// GO_ENV := os.Getenv("GO_ENV")
	// if GO_ENV != "prod" {
	// err := godotenv.Load()
	if err != nil {
		panic(`failed to load env: ` + err.Error())
	}
	// }

	cfg, err := config.LoadFromEnv()
	if err != nil {
		panic(`failed to init config: ` + err.Error())
	}

	logger, err := logging.NewZapLogger("prod")
	if err != nil {
		panic(`failed to init logger: ` + err.Error())
	}
	defer func() {
		_ = logger.Sync()
	}()

	botAPI, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		logger.Error("failed to init crypto-knight telegram bot", "error", err)
		os.Exit(1)
	}
	botAPI.Debug = false

	logger.Info("crypto-knight telegram bot started", "username", botAPI.Self.UserName)

	client := httpclient.New(cfg.APIBaseURL, time.Duration(cfg.HTTPTimeoutSeconds)*time.Second)

	appl := app.NewApp(botAPI, cfg, client)
	if err := appl.Run(ctx); err != nil {
		logger.Error("failed to run app", "error", err)
		os.Exit(1)
	}
}
