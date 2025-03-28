package main

import (
	"context"
	"os"
	"time"

	"google.golang.org/api/calendar/v3"

	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

type CalendarService interface {
	Init(ctx context.Context) error
	GetRecentEvents(ctx context.Context, since time.Time) ([]CalendarEvent, error)
}

type EventStatus int

const (
	StatusCreated EventStatus = iota
	StatusUpdated
	StatusCanceled
	StatusUnknown
)

type CalendarEvent struct {
	Title   string
	Start   time.Time
	End     time.Time
	Creator string
	Status  EventStatus
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
		ShowDeleted(true).
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
			Start:   parseEventStart(e),
			End:     parseEventEnd(e),
			Creator: getEventCreator(e),
			Status:  parseEventStatus(e),
		})
	}
	return resp, nil
}

const (
	googleStatusConfirmed = "confirmed"
	googleStatusCancelled = "cancelled"
)

func getEventCreator(event *calendar.Event) string {
	if event.Creator == nil {
		return ""
	}

	return event.Creator.Email
}

func parseEventStatus(event *calendar.Event) EventStatus {
	switch event.Status {
	case googleStatusConfirmed:
		if isEventUpdated(event) {
			return StatusUpdated
		} else {
			return StatusCreated
		}
	case googleStatusCancelled:
		return StatusCanceled
	default:
		return StatusUnknown
	}
}

func isEventUpdated(event *calendar.Event) bool {
	created, err := time.Parse(time.RFC3339, event.Created)
	if err != nil {
		return false
	}

	updated, err := time.Parse(time.RFC3339, event.Updated)
	if err != nil {
		return false
	}

	return updated.Sub(created) >= time.Second
}

func parseEventStart(event *calendar.Event) time.Time {
	return parseEventTime(event.Start)
}

func parseEventEnd(event *calendar.Event) time.Time {
	return parseEventTime(event.End)
}

func parseEventTime(eventDateTime *calendar.EventDateTime) time.Time {
	if eventDateTime.Date != "" {
		t, err := time.Parse(time.DateOnly, eventDateTime.Date)
		if err != nil {
			log.WithError(err).Error("error parsing event date")
			return time.Time{}
		}
		return t
	}

	t, err := time.Parse(time.RFC3339, eventDateTime.DateTime)
	if err != nil {
		log.WithError(err).Error("error parsing event time")
		return time.Time{}
	}

	return t
}
