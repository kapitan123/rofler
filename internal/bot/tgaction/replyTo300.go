package tgaction

import (
	"context"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/bot"
)

const matcher300 = `300|Триста|триста`

type ReplyTo300 struct {
	*bot.Bot
}

func NewReplyTo300(b *bot.Bot) *ReplyTo300 {
	return &ReplyTo300{
		Bot: b,
	}
}

func (h *ReplyTo300) Handle(m *tgbotapi.Message, ctx context.Context) (bool, error) {
	r := regexp.MustCompile(matcher300)
	found := r.MatchString(m.Text)

	if !found {
		return false, nil
	}

	h.PostReplyWithText(m.Chat.ID, m.MessageID, "🤣🚜 ♂ Отсоси у тракториста ♂ 🚜🤣")

	return false, nil
}
