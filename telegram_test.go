package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	name     string
	status   EventStatus
	expected string
}

func TestPrepareMessageBodyDateTime(t *testing.T) {
	title := "some title"
	start := time.Now().Add(24 * time.Hour)
	end := start.Add(time.Hour)

	testsDateTime := []testCase{
		{
			name:     "newly created event",
			status:   StatusCreated,
			expected: fmt.Sprintf("🗓️ *%s*\n\n*התחלה:* %s\n*סיום:* %s", title, FormatDateTime(start), FormatDateTime(end)),
		},
		{
			name:     "updated event",
			status:   StatusUpdated,
			expected: fmt.Sprintf("️✍🏻 *עדכון: %s*\n\n*התחלה:* %s\n*סיום:* %s", title, FormatDateTime(start), FormatDateTime(end)),
		},
		{
			name:     "cancelled event",
			status:   StatusCanceled,
			expected: fmt.Sprintf("️🆇 *בוטל: %s*\n\n*התחלה:* %s\n*סיום:* %s", title, FormatDateTime(start), FormatDateTime(end)),
		},
	}

	for _, test := range testsDateTime {
		t.Run(test.name, func(t *testing.T) {
			event := CalendarEvent{
				Title:  title,
				Start:  start,
				End:    end,
				Status: test.status,
			}
			actual, err := prepareMessageBody(event)
			require.NoError(t, err)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestPrepareMessageBodyDateOnly(t *testing.T) {
	title := "some title"
	start := time.Now().Add(24 * time.Hour)
	end := start.Add(time.Hour)

	startDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	endDay := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())

	testsDateOnly := []testCase{
		{
			name:     "newly created event",
			status:   StatusCreated,
			expected: fmt.Sprintf("🗓️ *%s*\n\n*התחלה:* %s\n*סיום:* %s", title, FormatDateTime(startDay), FormatDateTime(endDay)),
		},
		{
			name:     "updated event",
			status:   StatusUpdated,
			expected: fmt.Sprintf("️✍🏻 *עדכון: %s*\n\n*התחלה:* %s\n*סיום:* %s", title, FormatDateTime(startDay), FormatDateTime(endDay)),
		},
		{
			name:     "cancelled event",
			status:   StatusCanceled,
			expected: fmt.Sprintf("️🆇 *בוטל: %s*\n\n*התחלה:* %s\n*סיום:* %s", title, FormatDateTime(startDay), FormatDateTime(endDay)),
		},
	}

	for _, test := range testsDateOnly {
		t.Run(test.name, func(t *testing.T) {
			event := CalendarEvent{
				Title:  title,
				Start:  startDay,
				End:    endDay,
				Status: test.status,
			}
			actual, err := prepareMessageBody(event)
			require.NoError(t, err)
			assert.Equal(t, test.expected, actual)
		})
	}
}
