package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram interface {
	Init() error
	SendMessage(message string) error
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

func (t *telegram) SendMessage(message string) error {
	msg := tgbotapi.NewMessage(t.cfg.TelegramChatId, message)
	msg.ParseMode = "markdown"
	_, err := t.bot.Send(msg)
	return err
}
