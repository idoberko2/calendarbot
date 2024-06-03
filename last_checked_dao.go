package main

import (
	"io/fs"
	"os"
	"time"
)

type LastCheckedDao interface {
	GetLastChecked() (time.Time, bool, error)
	SetLastChecked(t time.Time) error
}

func NewLastCheckedDao(cfg Config) LastCheckedDao {
	return &lastCheckedDao{
		cfg: cfg,
	}
}

type lastCheckedDao struct {
	cfg Config
}

func (l *lastCheckedDao) GetLastChecked() (time.Time, bool, error) {
	// Read the time string from the file
	timeBytes, err := os.ReadFile(l.cfg.LastCheckedFile)
	if os.IsNotExist(err) {
		return time.Time{}, false, nil
	}
	if err != nil {
		return time.Time{}, false, err
	}

	// Convert the time string back to a time.Time object
	t, err := time.Parse(time.RFC3339, string(timeBytes))
	if err != nil {
		return time.Time{}, false, err
	}

	return t, true, nil
}

func (l *lastCheckedDao) SetLastChecked(t time.Time) error {
	// Convert the time to a string
	timeString := t.Format(time.RFC3339)

	// Write the time string to the file
	return os.WriteFile(l.cfg.LastCheckedFile, []byte(timeString), fs.FileMode(0644))
}
