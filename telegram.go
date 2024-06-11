package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
		return fmt.Sprintf("🗓️ *%s*\n\n*התחלה:* %s\n*סיום:* %s", event.Title, event.Start, event.End), nil
	case StatusUpdated:
		return fmt.Sprintf("️✍🏻 *עדכון: %s*\n\n*התחלה:* %s\n*סיום:* %s", event.Title, event.Start, event.End), nil
	case StatusCanceled:
		return fmt.Sprintf("️🆇 *בוטל: %s*\n\n*התחלה:* %s\n*סיום:* %s", event.Title, event.Start, event.End), nil
	default:
		return "", fmt.Errorf("unexpected status: %d", event.Status)
	}
}
