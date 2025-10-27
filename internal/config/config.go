package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	BotToken           string
	UserIDs            []int64
	NotificationGroup  int64
	APIBaseURL         string
	HTTPTimeoutSeconds int
	RmqURL             string
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

	return &Config{
		BotToken:           bot,
		UserIDs:            res,
		NotificationGroup:  groupID,
		APIBaseURL:         api,
		HTTPTimeoutSeconds: timeout,
		RmqURL:             rmqURL,
	}, nil
}
