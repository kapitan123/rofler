package tgaction

import (
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/bot"
)

var noRegex = regexp.MustCompile(`^(Net|net|Нет|нет)(.|\?|!)?$`)

type ReplyToNo struct {
	*bot.Bot
}

func NewReplyToNo(b *bot.Bot) *ReplyToYes {
	return &ReplyToYes{
		Bot: b,
	}
}

func (h *ReplyToNo) Handle(m *tgbotapi.Message) (bool, error) {
	// AK TODO post a picture with watermark to the chat
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
