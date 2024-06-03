package main

import (
	"context"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	if err := LoadDotEnv(); err != nil {
		log.WithError(err).Fatal("error loading .env file")
	}
	cfg, err := ReadConfigFromEnv(ctx)
	if err != nil {
		log.WithError(err).Fatal("error reading config")
	}

	calSvc := NewCalendarService(cfg)
	if err := calSvc.Init(ctx); err != nil {
		log.WithError(err).Fatal("error initializing calendar service")
	}

	telcli := NewTelegram(cfg)
	if err := telcli.Init(); err != nil {
		log.WithError(err).Fatal("error initializing telegram client")
	}

	lastChkdDao := NewLastCheckedDao(cfg)
	engine := Engine{
		cfg:         cfg,
		calSvc:      calSvc,
		telcli:      telcli,
		lastChkdDao: lastChkdDao,
	}

	if err := engine.Work(ctx); err != nil {
		log.WithError(err).Fatal("engine failed")
	}
}
