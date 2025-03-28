package main

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type EngineSuite struct {
	suite.Suite
	filename    string
	calSvcMock  *CalendarServiceMock
	telCliMock  *TelegramClientMock
	lastChkdDao LastCheckedDao
	calendarId  string
	engine      Engine
}

func TestEngineSuite(t *testing.T) {
	suite.Run(t, new(EngineSuite))
}

func (s *EngineSuite) SetupSuite() {
	s.filename = fmt.Sprintf("%s_%d.txt", "test_last_checked", time.Now().Unix())
	s.calendarId = "someCalendarId"
}

func (s *EngineSuite) SetupTest() {
	s.calSvcMock = &CalendarServiceMock{}
	s.telCliMock = &TelegramClientMock{}
	s.lastChkdDao = NewLastCheckedDao(Config{LastCheckedFile: s.filename})
	s.engine = Engine{
		cfg: Config{
			CalendarId: s.calendarId,
		},
		calSvc:      s.calSvcMock,
		telcli:      s.telCliMock,
		lastChkdDao: s.lastChkdDao,
	}

	_ = os.Remove(s.filename)
}

func (s *EngineSuite) TearDownSuite() {
	_ = os.Remove(s.filename)
}

func (s *EngineSuite) TestFirstRun() {
	ctx := context.Background()
	s.calSvcMock.On("GetRecentEvents", ctx, mock.Anything).Return([]CalendarEvent{}, nil)

	// SUT
	s.Require().NoError(s.engine.Work(ctx))

	s.Require().Len(s.calSvcMock.Calls, 1)
	calledTime, ok := s.calSvcMock.Calls[0].Arguments[1].(time.Time)
	s.Require().True(ok)
	s.Assert().WithinDuration(time.Now().Add(-time.Hour), calledTime, 5*time.Second)

	newLastChecked, isExist, err := s.lastChkdDao.GetLastChecked()
	s.Require().NoError(err)
	s.Assert().True(isExist)
	s.Assert().WithinDuration(time.Now(), newLastChecked, 5*time.Second)
}

func (s *EngineSuite) TestReceiveEvents() {
	ctx := context.Background()
	lastChecked := time.Now().Add(-time.Minute)
	s.Require().NoError(s.lastChkdDao.SetLastChecked(lastChecked))
	start := time.Now().Add(24 * time.Hour)
	end := start.Add(time.Hour)
	notifiedEvent := CalendarEvent{
		Title:   "Should be notified",
		Start:   start,
		End:     end,
		Creator: "someone else",
	}
	ignoredEvent := CalendarEvent{
		Title:   "Should be ignored",
		Start:   start,
		End:     end,
		Creator: s.calendarId,
	}
	s.calSvcMock.On("GetRecentEvents", ctx, mock.Anything).Return([]CalendarEvent{
		notifiedEvent,
		ignoredEvent,
	}, nil)
	s.telCliMock.On("NotifyEvent", mock.Anything).Return(nil)

	// SUT
	s.Require().NoError(s.engine.Work(ctx))

	s.Require().Len(s.calSvcMock.Calls, 1)
	calledTime, ok := s.calSvcMock.Calls[0].Arguments[1].(time.Time)
	s.Require().True(ok)
	s.Assert().WithinDuration(lastChecked, calledTime, time.Second)

	s.telCliMock.AssertNumberOfCalls(s.T(), "NotifyEvent", 1)
	s.telCliMock.AssertCalled(s.T(), "NotifyEvent", notifiedEvent)
}

func (s *EngineSuite) TestIgnoreEvents() {
	start := time.Now().Add(24 * time.Hour)
	end := start.Add(time.Hour)

	tests := []struct {
		name  string
		event CalendarEvent
	}{
		{"ignore configured calendar id", CalendarEvent{
			Title:   "Should be ignored",
			Start:   start,
			End:     end,
			Creator: s.calendarId,
		}},
		{"ignore outdated events", CalendarEvent{
			Title:   "Should be ignored",
			Start:   time.Now().Add(-24 * time.Hour),
			End:     time.Now().Add(-23 * time.Hour),
			Creator: "some-other-calendar-id",
		}},
	}
	ctx := context.Background()

	for _, test := range tests {
		s.Suite.Run(test.name, func() {
			lastChecked := time.Now().Add(-time.Minute)
			s.SetupTest()

			s.Require().NoError(s.lastChkdDao.SetLastChecked(lastChecked))
			s.calSvcMock.On("GetRecentEvents", ctx, mock.Anything).Return([]CalendarEvent{
				test.event,
			}, nil)
			s.telCliMock.On("NotifyEvent", mock.Anything).Return(nil)

			// SUT
			s.Require().NoError(s.engine.Work(ctx))

			s.Require().Len(s.calSvcMock.Calls, 1)
			calledTime, ok := s.calSvcMock.Calls[0].Arguments[1].(time.Time)
			s.Require().True(ok)
			s.Assert().WithinDuration(lastChecked, calledTime, time.Second)

			s.telCliMock.AssertNotCalled(s.T(), "NotifyEvent", mock.Anything)
		})
	}
}
