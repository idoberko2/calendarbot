package main

import (
	"context"
	"time"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

type Engine struct {
	cfg         Config
	calSvc      CalendarService
	telcli      Telegram
	lastChkdDao LastCheckedDao
}

func (e *Engine) Work(ctx context.Context) error {
	t, isExist, err := e.lastChkdDao.GetLastChecked()
	if !isExist {
		t = time.Now().Add(-time.Hour)
	} else if err != nil {
		return errors.Wrap(err, "error reading last checked from file")
	}

	timeCheck := time.Now()
	events, err := e.calSvc.GetRecentEvents(ctx, t)
	if err != nil {
		return errors.Wrap(err, "error getting events")
	}

	for _, event := range events {
		if event.Creator == e.cfg.CalendarId {
			log.Info("ignoring event created by calendar owner")
			continue
		}

		if event.Start.Before(timeCheck) {
			log.Info("ignoring outdated event")
			continue
		}

		if err := e.telcli.NotifyEvent(event); err != nil {
			return errors.Wrap(err, "error sending telegram message")
		}
	}

	if err := e.lastChkdDao.SetLastChecked(timeCheck); err != nil {
		return errors.Wrap(err, "error writing last checked time")
	}

	return nil
}
