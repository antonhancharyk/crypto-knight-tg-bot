package telegram

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/config"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type userFlowState struct {
	mu   sync.Mutex
	From string
	To   string
	Step int // 0 none, 1 waiting from, 2 waiting to
}

type Handler struct {
	bot      *tgbotapi.BotAPI
	cfg      *config.Config
	reportUC *usecase.ReportUsecase
	states   map[int64]*userFlowState
	statesMu sync.Mutex
}

func NewHandler(bot *tgbotapi.BotAPI, cfg *config.Config, ru *usecase.ReportUsecase) *Handler {
	return &Handler{bot: bot, cfg: cfg, reportUC: ru, states: make(map[int64]*userFlowState)}
}

func (h *Handler) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := h.bot.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			h.bot.StopReceivingUpdates()
			return nil

		case update, ok := <-updates:
			if !ok {
				return nil
			}
			if update.Message != nil {
				go h.handleMessage(update.Message)
				continue
			}
			if update.CallbackQuery != nil {
				go h.handleCallback(update.CallbackQuery)
			}
		}
	}
}

func (h *Handler) isAdmin(id int64) bool {
	return slices.Contains(h.cfg.UserIDs, id)
}

func (h *Handler) getState(userID int64) *userFlowState {
	h.statesMu.Lock()
	defer h.statesMu.Unlock()
	st, ok := h.states[userID]
	if !ok {
		st = &userFlowState{}
		h.states[userID] = st
	}
	return st
}

func (h *Handler) handleMessage(msg *tgbotapi.Message) {
	userID := msg.From.ID
	chatID := msg.Chat.ID

	if !h.isAdmin(userID) {
		_ = h.reply(chatID, "Access denied")
		return
	}

	text := strings.TrimSpace(msg.Text)
	st := h.getState(userID)
	st.mu.Lock()
	defer st.mu.Unlock()

	if text == "/start" {
		h.sendMenu(chatID)
		st.Step = 0
		return
	}

	// flow logic
	switch st.Step {
	case 1:
		st.From = text
		st.Step = 2
		_ = h.reply(chatID, "Enter end date (YYYY-MM-DD):")
		return
	case 2:
		st.To = text
		st.Step = 0
		// call usecase
		ctx := context.Background()
		rep, err := h.reportUC.GetReport(ctx, st.From, st.To)
		if err != nil {
			_ = h.reply(chatID, fmt.Sprintf("Error: %v", err))
			return
		}
		_ = h.reply(chatID, fmt.Sprintf("Report from %s to %s. Income: %.2f, Expense: %.2f", rep.From, rep.To, rep.Income, rep.Expense))
		return
	}

	_ = h.reply(chatID, "Unknown input. Use /start to open menu")
}

func (h *Handler) handleCallback(q *tgbotapi.CallbackQuery) {
	data := q.Data
	userID := q.From.ID
	chatID := q.Message.Chat.ID

	if !h.isAdmin(userID) {
		_ = h.answerCallback(q, "Access denied")
		return
	}

	switch data {
	case "menu:report":
		// start report flow
		st := h.getState(userID)
		st.mu.Lock()
		st.Step = 1
		st.From = ""
		st.To = ""
		st.mu.Unlock()
		_ = h.answerCallback(q, "Enter start date (YYYY-MM-DD):")
		_ = h.reply(chatID, "Enter start date (YYYY-MM-DD):")
	default:
		_ = h.answerCallback(q, "Unknown action")
	}
}

func (h *Handler) sendMenu(chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, "Choose action:")
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Get income/expenses", "menu:report"),
		),
	)
	msg.ReplyMarkup = kb
	_, err := h.bot.Send(msg)
	return err
}

func (h *Handler) reply(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := h.bot.Send(msg)
	return err
}

func (h *Handler) answerCallback(q *tgbotapi.CallbackQuery, text string) error {
	cfg := tgbotapi.NewCallback(q.ID, text)
	_, err := h.bot.Request(cfg)
	return err
}
