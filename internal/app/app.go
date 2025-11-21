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

	consumers := []struct {
		queue   string
		handler func([]byte) error
	}{
		{"trading-signals-queue", func(msg []byte) error {
			return h.SendToGroup(-4603798918, string(msg))
		}},
		{"pnl-reports-queue", func(msg []byte) error {
			return h.SendToGroup(-5082938682, string(msg))
		}},
		{"system-queue", func(msg []byte) error {
			return h.SendToGroup(-1003283451332, string(msg))
		}},
	}

	for _, cfg := range consumers {
		consumer := broker.NewConsumer(a.rmq.Channel(), cfg.queue, cfg.handler)
		if err := consumer.Run(ctx); err != nil {
			return err
		}
	}

	h.Run(ctx)

	return nil
}
