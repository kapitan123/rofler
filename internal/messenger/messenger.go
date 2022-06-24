package messenger

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Messenger struct {
	api *tgbotapi.BotAPI
}

func New(api *tgbotapi.BotAPI) *Messenger {
	return &Messenger{api: api}
}

func (m *Messenger) SendMessage(_ context.Context, chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := m.api.Send(msg)
	return err
}
