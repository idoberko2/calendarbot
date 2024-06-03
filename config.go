package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-envconfig"
	log "github.com/sirupsen/logrus"
	"io/fs"
)

type Config struct {
	CalendarId               string `env:"CALENDAR_ID"`
	TelegramToken            string `env:"TELEGRAM_TOKEN"`
	TelegramChatId           int64  `env:"TELEGRAM_CHAT_ID"`
	LastCheckedFile          string `env:"LAST_CHECKED_FILE, default=last_checked.txt"`
	GoogleServiceAccountFile string `env:"GOOGLE_SERVICE_ACCOUNT_FILE, default=google_service_account.json"`
}

func ReadConfigFromEnv(ctx context.Context) (Config, error) {
	var cfg Config

	if err := envconfig.Process(ctx, &cfg); err != nil {
		return cfg, errors.Wrap(err, "error processing config")
	}

	return cfg, nil
}

func LoadDotEnv() error {
	var pathErr *fs.PathError

	if err := godotenv.Load(".env"); errors.As(err, &pathErr) {
		log.Info("couldn't find .env file, skipping .env file load")
	} else if err != nil {
		return err
	}

	return nil
}
