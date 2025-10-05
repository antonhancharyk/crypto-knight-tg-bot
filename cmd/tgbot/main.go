package main

import (
	"log"
	"time"

	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/app"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/config"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/infra/httpclient"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	botAPI, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatalf("telegram init: %v", err)
	}
	botAPI.Debug = false

	client := httpclient.New(cfg.APIBaseURL, time.Duration(cfg.HTTPTimeoutSeconds)*time.Second)

	appl := app.NewApp(botAPI, cfg, client)
	if err := appl.Run(); err != nil {
		log.Fatalf("app run: %v", err)
	}
}
