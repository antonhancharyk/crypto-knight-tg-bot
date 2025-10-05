package app

import (
	"context"
	"log"

	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/config"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/infra/httpclient"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/transport/telegram"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	botAPI *tgbotapi.BotAPI
	cfg    *config.Config
	client *httpclient.Client
}

func NewApp(botAPI *tgbotapi.BotAPI, cfg *config.Config, client *httpclient.Client) *App {
	return &App{botAPI: botAPI, cfg: cfg, client: client}
}

func (a *App) Run() error {
	ruc := usecase.NewReportUsecase(a.client)
	h := telegram.NewHandler(a.botAPI, a.cfg, ruc)
	log.Println("tgbot started")
	return h.Run(context.Background())
}
