package main

import (
	"context"
	"fmt"
	"time"

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
	t, isExist, err := lastChkdDao.GetLastChecked()
	if !isExist {
		t = time.Now().Add(-time.Hour)
	} else if err != nil {
		log.WithError(err).Fatal("error reading last checked from file")
	}

	timeCheck := time.Now()
	events, err := calSvc.GetRecentEvents(ctx, t)
	if err != nil {
		log.WithError(err).Fatal("error getting events")
	}

	for _, e := range events {
		if e.Creator == cfg.CalendarId {
			continue
		}

		msg := fmt.Sprintf("🗓️ *%s*\n\n*התחלה:* %s\n*סיום:* %s", e.Title, e.Start, e.End)
		if err := telcli.SendMessage(msg); err != nil {
			log.WithError(err).Fatal("error sending telegram message")
		}
	}

	if err := lastChkdDao.SetLastChecked(timeCheck); err != nil {
		log.WithError(err).Fatal("error writing last checked time")
	}
}
