package main

import (
	"context"
	"google.golang.org/api/calendar/v3"
	"os"
	"time"

	"google.golang.org/api/option"
)

type CalendarService interface {
	Init(ctx context.Context) error
	GetRecentEvents(ctx context.Context, since time.Time) ([]CalendarEvent, error)
}

type CalendarEvent struct {
	Title   string
	Start   string
	End     string
	Creator string
}

func NewCalendarService(cfg Config) CalendarService {
	return &calendarClient{
		cfg: cfg,
	}
}

type calendarClient struct {
	cfg Config
	svc *calendar.Service
}

func (c *calendarClient) Init(ctx context.Context) error {
	b, err := os.ReadFile(c.cfg.GoogleServiceAccountFile)
	if err != nil {
		return err
	}

	svc, err := calendar.NewService(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		return err
	}

	c.svc = svc
	return nil
}

func (c *calendarClient) GetRecentEvents(ctx context.Context, since time.Time) ([]CalendarEvent, error) {
	events, err := c.svc.Events.
		List(c.cfg.CalendarId).
		ShowDeleted(false).
		SingleEvents(true).
		UpdatedMin(since.Format(time.RFC3339)).
		MaxResults(10).
		OrderBy("updated").
		Context(ctx).
		Do()
	if err != nil {
		return nil, err
	}

	resp := make([]CalendarEvent, 0, len(events.Items))
	for _, e := range events.Items {
		resp = append(resp, CalendarEvent{
			Title:   e.Summary,
			Start:   e.Start.DateTime,
			End:     e.End.DateTime,
			Creator: e.Creator.Email,
		})
	}
	return resp, nil
}
