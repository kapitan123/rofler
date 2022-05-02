package tgaction

import (
	_ "embed"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/bot"
)

var yesRegex = regexp.MustCompile(`^(da|Да|да)(.|\?|!)?$`)

//go:embed assets/kirkorov.png
var yesPicture []byte

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

	err := h.PostReplyWithImage(m.Chat.ID, m.MessageID, yesPicture, "kirkorov.png", "")

	if err != nil {
		return false, err
	}
	return true, nil
}
