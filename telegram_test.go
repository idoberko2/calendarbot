package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrepareMessageBody(t *testing.T) {
	title := "some title"
	start := time.Now().Add(24 * time.Hour)
	end := start.Add(time.Hour)
	tests := []struct {
		name     string
		status   EventStatus
		expected string
	}{
		{
			name:     "newly created event",
			status:   StatusCreated,
			expected: fmt.Sprintf("🗓️ *%s*\n\n*התחלה:* %s\n*סיום:* %s", title, start, end),
		},
		{
			name:     "updated event",
			status:   StatusUpdated,
			expected: fmt.Sprintf("️✍🏻 *עדכון: %s*\n\n*התחלה:* %s\n*סיום:* %s", title, start, end),
		},
		{
			name:     "cancelled event",
			status:   StatusCanceled,
			expected: fmt.Sprintf("️🆇 *בוטל: %s*\n\n*התחלה:* %s\n*סיום:* %s", title, start, end),
		},
	}
	for _, test := range tests {
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
