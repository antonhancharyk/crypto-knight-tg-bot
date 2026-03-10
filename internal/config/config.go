package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

// QueueConsumer defines a RabbitMQ queue and the Telegram group to forward messages to.
type QueueConsumer struct {
	QueueName   string
	GroupChatID int64
}

type Config struct {
	BotToken           string
	UserIDs            []int64
	NotificationGroup  int64
	APIBaseURL         string
	HTTPTimeoutSeconds int
	RmqURL             string
	// QueueConsumers defines queues to consume and their target group chat IDs.
	QueueConsumers []QueueConsumer
}

func LoadFromEnv() (*Config, error) {
	bot := os.Getenv("BOT_TOKEN")
	if bot == "" {
		return nil, errors.New("BOT_TOKEN required")
	}

	users := os.Getenv("ADMIN_USER_IDS")
	if users == "" {
		return nil, errors.New("ADMIN_USER_IDS required")
	}

	api := os.Getenv("API_BASE_URL")
	if api == "" {
		api = "http://localhost:8081"
	}

	timeout := 10
	if s := os.Getenv("HTTP_TIMEOUT"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			timeout = v
		}
	}

	var groupID int64 = 0
	if g := os.Getenv("NOTIFICATION_GROUP_ID"); g != "" {
		if v, err := strconv.ParseInt(g, 10, 64); err == nil {
			groupID = v
		}
	}

	parts := strings.Split(users, ",")
	res := make([]int64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if v, err := strconv.ParseInt(p, 10, 64); err == nil {
			res = append(res, v)
		}
	}

	rmqURL := os.Getenv("RABBITMQ_URL")
	if rmqURL == "" {
		return nil, errors.New("RABBITMQ_URL required")
	}

	queueConsumers := loadQueueConsumers()

	return &Config{
		BotToken:           bot,
		UserIDs:            res,
		NotificationGroup:  groupID,
		APIBaseURL:         api,
		HTTPTimeoutSeconds: timeout,
		RmqURL:             rmqURL,
		QueueConsumers:     queueConsumers,
	}, nil
}

// loadQueueConsumers reads queue/group pairs from env with defaults for backward compatibility.
func loadQueueConsumers() []QueueConsumer {
	// Optional: CONSUMER_QUEUES="queue1:chatId1,queue2:chatId2" (comma-separated queue:chatId)
	// If unset, use default queues and env-based group IDs.
	defaults := []struct {
		queueEnv, groupEnv string
		queueDefault       string
		groupDefault       int64
	}{
		{"TRADING_SIGNALS_QUEUE", "TRADING_SIGNALS_GROUP_ID", "trading-signals-queue", -4603798918},
		{"PNL_REPORTS_QUEUE", "PNL_REPORTS_GROUP_ID", "pnl-reports-queue", -5082938682},
		{"SYSTEM_QUEUE", "SYSTEM_GROUP_ID", "system-queue", -1003283451332},
	}
	out := make([]QueueConsumer, 0, len(defaults))
	for _, d := range defaults {
		q := os.Getenv(d.queueEnv)
		if q == "" {
			q = d.queueDefault
		}
		g := d.groupDefault
		if s := os.Getenv(d.groupEnv); s != "" {
			if v, err := strconv.ParseInt(s, 10, 64); err == nil {
				g = v
			}
		}
		out = append(out, QueueConsumer{QueueName: q, GroupChatID: g})
	}
	return out
}
