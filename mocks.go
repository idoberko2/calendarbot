package main

import (
	"context"
	"github.com/stretchr/testify/mock"
	"time"
)

type CalendarServiceMock struct {
	mock.Mock
}

func (c *CalendarServiceMock) Init(ctx context.Context) error {
	args := c.Called(ctx)
	return args.Error(0)
}

func (c *CalendarServiceMock) GetRecentEvents(ctx context.Context, since time.Time) ([]CalendarEvent, error) {
	args := c.Called(ctx, since)
	return args.Get(0).([]CalendarEvent), args.Error(1)
}

type TelegramClientMock struct {
	mock.Mock
}

func (t *TelegramClientMock) Init() error {
	args := t.Called()
	return args.Error(0)
}

func (t *TelegramClientMock) NotifyEvent(event CalendarEvent) error {
	args := t.Called(event)
	return args.Error(0)
}
