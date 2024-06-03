package main

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

type LastCheckedDaoSuite struct {
	suite.Suite
	filename string
	dao      LastCheckedDao
}

func TestLastCheckedDaoSuite(t *testing.T) {
	suite.Run(t, new(LastCheckedDaoSuite))
}

func (s *LastCheckedDaoSuite) SetupSuite() {
	s.filename = fmt.Sprintf("%s_%d.txt", "test_last_checked", time.Now().Unix())
	s.dao = NewLastCheckedDao(Config{LastCheckedFile: s.filename})
}

func (s *LastCheckedDaoSuite) SetupTest() {
	_ = os.Remove(s.filename)
}

func (s *LastCheckedDaoSuite) TearDownSuite() {
	_ = os.Remove(s.filename)
}

func (s *LastCheckedDaoSuite) TestGetBeforeSet() {
	_, isExist, err := s.dao.GetLastChecked()
	s.Assert().False(isExist)
	s.Assert().NoError(err)
}

func (s *LastCheckedDaoSuite) TestSetGet() {
	expected, err := time.Parse(time.RFC3339, "2024-06-03T21:28:00Z")
	s.Require().NoError(err)
	s.Require().NoError(s.dao.SetLastChecked(expected))

	actual, isExist, err := s.dao.GetLastChecked()
	s.Assert().NoError(err)
	s.Assert().True(isExist)
	s.Assert().Equal(expected, actual)
}
