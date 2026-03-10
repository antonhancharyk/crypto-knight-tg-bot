package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadFromEnv(t *testing.T) {
	// Save and restore env to avoid polluting the process
	save := map[string]string{}
	for _, k := range []string{"BOT_TOKEN", "ADMIN_USER_IDS", "RABBITMQ_URL", "API_BASE_URL", "HTTP_TIMEOUT", "NOTIFICATION_GROUP_ID"} {
		save[k] = os.Getenv(k)
	}
	t.Cleanup(func() {
		for k, v := range save {
			if v == "" {
				os.Unsetenv(k)
			} else {
				_ = os.Setenv(k, v)
			}
		}
	})

	t.Run("missing BOT_TOKEN", func(t *testing.T) {
		os.Unsetenv("BOT_TOKEN")
		_ = os.Setenv("ADMIN_USER_IDS", "123")
		_ = os.Setenv("RABBITMQ_URL", "amqp://localhost")
		_, err := LoadFromEnv()
		require.Error(t, err)
		require.Contains(t, err.Error(), "BOT_TOKEN")
	})

	t.Run("missing ADMIN_USER_IDS", func(t *testing.T) {
		_ = os.Setenv("BOT_TOKEN", "token")
		os.Unsetenv("ADMIN_USER_IDS")
		_ = os.Setenv("RABBITMQ_URL", "amqp://localhost")
		_, err := LoadFromEnv()
		require.Error(t, err)
		require.Contains(t, err.Error(), "ADMIN_USER_IDS")
	})

	t.Run("missing RABBITMQ_URL", func(t *testing.T) {
		_ = os.Setenv("BOT_TOKEN", "token")
		_ = os.Setenv("ADMIN_USER_IDS", "123,456")
		os.Unsetenv("RABBITMQ_URL")
		_, err := LoadFromEnv()
		require.Error(t, err)
		require.Contains(t, err.Error(), "RABBITMQ_URL")
	})

	t.Run("success with defaults", func(t *testing.T) {
		_ = os.Setenv("BOT_TOKEN", "token")
		_ = os.Setenv("ADMIN_USER_IDS", "123, 456 ")
		_ = os.Setenv("RABBITMQ_URL", "amqp://guest:guest@localhost/")
		os.Unsetenv("API_BASE_URL")
		os.Unsetenv("HTTP_TIMEOUT")
		os.Unsetenv("NOTIFICATION_GROUP_ID")

		cfg, err := LoadFromEnv()
		require.NoError(t, err)
		require.Equal(t, "token", cfg.BotToken)
		require.Equal(t, []int64{123, 456}, cfg.UserIDs)
		require.Equal(t, "http://localhost:8081", cfg.APIBaseURL)
		require.Equal(t, 10, cfg.HTTPTimeoutSeconds)
		require.Equal(t, int64(0), cfg.NotificationGroup)
		require.Len(t, cfg.QueueConsumers, 3)
		require.Equal(t, "trading-signals-queue", cfg.QueueConsumers[0].QueueName)
		require.Equal(t, "pnl-reports-queue", cfg.QueueConsumers[1].QueueName)
		require.Equal(t, "system-queue", cfg.QueueConsumers[2].QueueName)
	})

	t.Run("optional env overrides", func(t *testing.T) {
		_ = os.Setenv("BOT_TOKEN", "t")
		_ = os.Setenv("ADMIN_USER_IDS", "1")
		_ = os.Setenv("RABBITMQ_URL", "amqp://x")
		_ = os.Setenv("API_BASE_URL", "https://api.example.com")
		_ = os.Setenv("HTTP_TIMEOUT", "30")
		_ = os.Setenv("NOTIFICATION_GROUP_ID", "-999")

		cfg, err := LoadFromEnv()
		require.NoError(t, err)
		require.Equal(t, "https://api.example.com", cfg.APIBaseURL)
		require.Equal(t, 30, cfg.HTTPTimeoutSeconds)
		require.Equal(t, int64(-999), cfg.NotificationGroup)
	})
}
