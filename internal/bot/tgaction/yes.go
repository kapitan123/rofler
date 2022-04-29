package tgaction

import (
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/bot"
)

var yesRegex = regexp.MustCompile(`^da|Да|да$`)

type ReplyToYes struct {
	*bot.Bot
}

func NewReplyToYes(b *bot.Bot) *ReplyToYes {
	return &ReplyToYes{
		Bot: b,
	}
}

func (h *ReplyToYes) Handle(m *tgbotapi.Message) (bool, error) {
	found := yesRegex.MatchString(m.Text)

	if !found {
		return false, nil
	}

	err := h.PostReplyToYes(m.Chat.ID, m.MessageID)
	if err != nil {
		return false, err
	}
	return true, nil
}
