package replyTo300

import (
	"context"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var regexp300 = regexp.MustCompile(`300|Триста|триста`)

type ReplyTo300 struct {
	messenger messenger
}

type messenger interface {
	ReplyWithText(chatId int64, messageId int, text string) error
}

func New(messenger messenger) *ReplyTo300 {
	return &ReplyTo300{
		messenger: messenger,
	}
}

func (h *ReplyTo300) Handle(ctx context.Context, m *tgbotapi.Message) error {
	err := h.messenger.ReplyWithText(m.Chat.ID, m.MessageID, "🤣🚜 ♂ Отсоси у тракториста ♂ 🚜🤣")

	return err
}

func (h *ReplyTo300) ShouldRun(m *tgbotapi.Message) bool {
	return regexp300.MatchString(m.Text)
}
