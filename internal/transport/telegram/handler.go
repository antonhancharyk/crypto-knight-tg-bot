package telegram

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/config"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type userFlowState struct {
	mu    sync.Mutex
	From  string
	To    string
	Step  int // 0 none, 1 waiting from, 2 waiting to
	Year  int
	Month int
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

	st := h.getState(userID)
	st.mu.Lock()
	defer st.mu.Unlock()

	if data == "menu:total_profit_loss" {
		st.Step = 1
		st.From = ""
		st.To = ""
		today := time.Now()
		_ = h.sendCalendar(chatID, 1, today.Year(), today.Month())
		return
	}

	if strings.HasPrefix(data, "date:") {
		parts := strings.Split(data, ":")
		date := parts[1]
		step := parts[2]

		if step == "1" {
			st.From = date
			st.Step = 2
			_ = h.answerCallback(q, "Start date selected: "+date)
			today := time.Now()
			_ = h.sendCalendar(chatID, 2, today.Year(), today.Month())
		} else if step == "2" {
			st.To = date
			st.Step = 0
			_ = h.answerCallback(q, "End date selected: "+date)

			ctx := context.Background()
			rep, err := h.reportUC.GetReport(ctx, st.From, st.To)
			if err != nil {
				_ = h.reply(chatID, fmt.Sprintf("error: %v", err))
				return
			}

			_ = h.reply(chatID, fmt.Sprintf("Report from %s to %s. Income: %.2f, Expense: %.2f", rep.From, rep.To, rep.Income, rep.Expense))
		}
		return
	}

	if strings.HasPrefix(data, "month:") {
		parts := strings.Split(data, ":")
		yearMonth := strings.Split(parts[1], "-")
		y, err := strconv.Atoi(yearMonth[0])
		if err != nil {
			_ = h.reply(chatID, fmt.Sprintf("convert error: %v", err))
		}
		m, err := strconv.Atoi(yearMonth[1])
		if err != nil {
			_ = h.reply(chatID, fmt.Sprintf("convert error: %v", err))
		}
		step := parts[3]

		var newMonth time.Month
		var newYear int
		t := time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.UTC)
		newYear, newMonth = t.Year(), t.Month()

		_ = h.answerCallback(q, "Month changed")
		i, err := strconv.Atoi(step)
		if err != nil {
			_ = h.reply(chatID, fmt.Sprintf("convert error: %v", err))
		}
		_ = h.sendCalendar(chatID, i, newYear, newMonth)
		return
	}

	_ = h.answerCallback(q, "Unknown action")
}

func (h *Handler) sendMenu(chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, "Choose action:")
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Total Profit/Loss %", "menu:total_profit_loss"),
		),
	)
	msg.ReplyMarkup = kb
	_, err := h.bot.Send(msg)
	return err
}

func (h *Handler) sendCalendar(chatID int64, step int, year int, month time.Month) error {
	daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()

	var rows [][]tgbotapi.InlineKeyboardButton
	week := []tgbotapi.InlineKeyboardButton{}

	for d := 1; d <= daysInMonth; d++ {
		dateStr := fmt.Sprintf("%04d-%02d-%02d", year, month, d)
		btn := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%02d", d), fmt.Sprintf("date:%s:%d", dateStr, step))
		week = append(week, btn)
		if len(week) == 7 {
			rows = append(rows, week)
			week = []tgbotapi.InlineKeyboardButton{}
		}
	}
	if len(week) > 0 {
		rows = append(rows, week)
	}

	prevMonth := time.Date(year, month-1, 1, 0, 0, 0, 0, time.UTC)
	nextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("◀", fmt.Sprintf("month:%d-%d:prev:%d", prevMonth.Year(), prevMonth.Month(), step)),
		tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s %d", month.String(), year), "noop"),
		tgbotapi.NewInlineKeyboardButtonData("▶", fmt.Sprintf("month:%d-%d:next:%d", nextMonth.Year(), nextMonth.Month(), step)),
	))

	kb := tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg := tgbotapi.NewMessage(chatID, "Select date:")
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
