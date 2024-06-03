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
	msgBody := fmt.Sprintf("🗓️ *%s*\n\n*התחלה:* %s\n*סיום:* %s", event.Title, event.Start, event.End)
	msg := tgbotapi.NewMessage(t.cfg.TelegramChatId, msgBody)
	msg.ParseMode = "markdown"
	_, err := t.bot.Send(msg)
	return err
}
