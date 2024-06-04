package main

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/calendar/v3"
	"testing"
	"time"
)

func TestParseStatus(t *testing.T) {
	baseTime := time.Date(2018, time.August, 19, 1, 2, 3, 4, time.UTC)

	tests := []struct {
		name     string
		input    *calendar.Event
		expected EventStatus
	}{
		{
			name: "newly created event - no time difference",
			input: &calendar.Event{
				Status:  googleStatusConfirmed,
				Created: baseTime.Format(time.RFC3339),
				Updated: baseTime.Format(time.RFC3339),
			},
			expected: StatusCreated,
		},
		{
			name: "newly created event - slight time difference",
			input: &calendar.Event{
				Status:  googleStatusConfirmed,
				Created: baseTime.Format(time.RFC3339),
				Updated: baseTime.Add(time.Millisecond).Format(time.RFC3339),
			},
			expected: StatusCreated,
		},
		{
			name: "updated event",
			input: &calendar.Event{
				Status:  googleStatusConfirmed,
				Created: baseTime.Format(time.RFC3339),
				Updated: baseTime.Add(time.Minute).Format(time.RFC3339),
			},
			expected: StatusUpdated,
		},
		{
			name: "cancelled event",
			input: &calendar.Event{
				Status:  googleStatusCancelled,
				Created: baseTime.Format(time.RFC3339),
				Updated: baseTime.Add(time.Minute).Format(time.RFC3339),
			},
			expected: StatusCanceled,
		},
		{
			name: "unknown event status",
			input: &calendar.Event{
				Status:  "unexpected google status",
				Created: baseTime.Format(time.RFC3339),
				Updated: baseTime.Add(time.Minute).Format(time.RFC3339),
			},
			expected: StatusUnknown,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := parseEventStatus(test.input)
			assert.Equal(t, test.expected, actual)
		})
	}
}
