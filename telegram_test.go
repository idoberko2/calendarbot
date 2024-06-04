package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrepareMessageBody(t *testing.T) {
	title := "some title"
	start := "start time"
	end := "end time"
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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			event := CalendarEvent{
				Title:  title,
				Start:  start,
				End:    end,
				Status: test.status,
			}
			actual := prepareMessageBody(event)
			assert.Equal(t, test.expected, actual)
		})
	}
}
