package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/calendar/v3"
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

func TestParseEventDates(t *testing.T) {
	tests := []struct {
		name          string
		input         *calendar.Event
		expectedStart time.Time
		expectedEnd   time.Time
	}{
		{
			name: "date and time",
			input: &calendar.Event{
				Start: &calendar.EventDateTime{
					DateTime: "2024-06-12T20:00:00+03:00",
				},
				End: &calendar.EventDateTime{
					DateTime: "2024-06-12T21:00:00+03:00",
				},
			},
			expectedStart: parseTimeNoError("2024-06-12T20:00:00+03:00"),
			expectedEnd:   parseTimeNoError("2024-06-12T21:00:00+03:00"),
		},
		{
			name: "date only",
			input: &calendar.Event{
				Start: &calendar.EventDateTime{
					Date: "2024-06-12",
				},
				End: &calendar.EventDateTime{
					Date: "2024-06-13",
				},
			},
			expectedStart: parseDateNoError("2024-06-12"),
			expectedEnd:   parseDateNoError("2024-06-12"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualStart := parseEventStart(test.input)
			actualEnd := parseEventEnd(test.input)
			assert.Equal(t, test.expectedStart, actualStart)
			assert.Equal(t, test.expectedEnd, actualEnd)
		})
	}
}

func parseTimeNoError(timeString string) time.Time {
	t, _ := time.Parse(time.RFC3339, timeString)
	return t
}

func parseDateNoError(dateString string) time.Time {
	t, _ := time.Parse(time.DateOnly, dateString)
	return t
}
