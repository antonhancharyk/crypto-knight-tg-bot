package app

import (
	"context"

	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/config"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/infra/broker"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/infra/httpclient"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/transport/telegram"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	botAPI *tgbotapi.BotAPI
	cfg    *config.Config
	client *httpclient.Client
	rmq    *broker.Connection
}

func NewApp(botAPI *tgbotapi.BotAPI, cfg *config.Config, client *httpclient.Client, rmq *broker.Connection) *App {
	return &App{botAPI: botAPI, cfg: cfg, client: client, rmq: rmq}
}

func (a *App) Run(ctx context.Context) error {
	ruc := usecase.NewReportUsecase(a.client)
	h := telegram.NewHandler(a.botAPI, a.cfg, ruc)

	consumer := broker.NewConsumer(a.rmq.Channel(), "tg", func(msg []byte) error {
		return h.SendToGroup(string(msg))
	})
	err := consumer.Run(ctx)
	if err != nil {
		return err
	}

	h.Run(ctx)

	return nil
}
