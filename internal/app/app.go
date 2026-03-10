package app

import (
	"context"

	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/config"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/infra/broker"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/ports"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/transport/telegram"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	botAPI  *tgbotapi.BotAPI
	cfg     *config.Config
	fetcher ports.ReportFetcher
	rmq     *broker.Connection
}

func NewApp(botAPI *tgbotapi.BotAPI, cfg *config.Config, fetcher ports.ReportFetcher, rmq *broker.Connection) *App {
	return &App{botAPI: botAPI, cfg: cfg, fetcher: fetcher, rmq: rmq}
}

func (a *App) Run(ctx context.Context) error {
	ruc := usecase.NewReportUsecase(a.fetcher)
	h := telegram.NewHandler(a.botAPI, a.cfg, ruc)

	for _, qc := range a.cfg.QueueConsumers {
		chatID := qc.GroupChatID
		consumer := broker.NewConsumer(a.rmq.Channel(), qc.QueueName, func(msg []byte) error {
			return h.SendToGroup(chatID, string(msg))
		})
		if err := consumer.Run(ctx); err != nil {
			return err
		}
	}

	h.Run(ctx)

	<-ctx.Done()
	return ctx.Err()
}
