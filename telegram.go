package main

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var weekdayDict = map[time.Weekday]string{
	time.Sunday:    "ראשון",
	time.Monday:    "שני",
	time.Tuesday:   "שלישי",
	time.Wednesday: "רביעי",
	time.Thursday:  "חמישי",
	time.Friday:    "שישי",
	time.Saturday:  "שבת",
}

type Telegram interface {
	Init() error
	NotifyEvent(event CalendarEvent) error
}

func NewTelegram(cfg Config) Telegram {
	return &telegram{
		cfg: cfg,
	}
}

type telegram struct {
	cfg Config
	bot *tgbotapi.BotAPI
}

func (t *telegram) Init() error {
	bot, err := tgbotapi.NewBotAPI(t.cfg.TelegramToken)
	if err != nil {
		return err
	}
	t.bot = bot

	return nil
}

func (t *telegram) NotifyEvent(event CalendarEvent) error {
	msgBody, err := prepareMessageBody(event)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(t.cfg.TelegramChatId, msgBody)
	msg.ParseMode = "markdown"
	_, err = t.bot.Send(msg)
	return err
}

func prepareMessageBody(event CalendarEvent) (string, error) {
	switch event.Status {
	case StatusCreated:
		return fmt.Sprintf(
			"🗓️ *%s*\n\n*התחלה:* %s\n*סיום:* %s",
			event.Title,
			FormatDateTime(event.Start),
			FormatDateTime(event.End)), nil
	case StatusUpdated:
		return fmt.Sprintf(
			"️✍🏻 *עדכון: %s*\n\n*התחלה:* %s\n*סיום:* %s",
			event.Title,
			FormatDateTime(event.Start),
			FormatDateTime(event.End)), nil
	case StatusCanceled:
		return fmt.Sprintf(
			"️🆇 *בוטל: %s*\n\n*התחלה:* %s\n*סיום:* %s",
			event.Title,
			FormatDateTime(event.Start),
			FormatDateTime(event.End)), nil
	default:
		return "", fmt.Errorf("unexpected status: %d", event.Status)
	}
}

func FormatDateTime(t time.Time) string {
	if t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 {
		return fmt.Sprintf("%s (%s)", t.Format(time.DateOnly), getDayOfWeek(t))
	}

	return fmt.Sprintf("%s (%s)", t.Format(time.DateTime), getDayOfWeek(t))
}

func getDayOfWeek(t time.Time) string {
	return weekdayDict[t.Weekday()]
}
