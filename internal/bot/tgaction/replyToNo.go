package tgaction

import (
	_ "embed"
	"regexp"

	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/pkg/watermarker"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var noRegex = regexp.MustCompile(`^(Net|net|Нет|нет)(.|\?|!)?$`)

type ReplyToNo struct {
	*bot.Bot
}

func NewReplyToNo(b *bot.Bot) *ReplyToNo {
	return &ReplyToNo{
		Bot: b,
	}
}

//go:embed assets/pidormark.png
var pidormarkPicture []byte

func (h *ReplyToNo) Handle(m *tgbotapi.Message) (bool, error) {
	// AK TODO post a picture with watermark to the chat
	found := noRegex.MatchString(m.Text)

	if !found {
		return false, nil
	}

	ppic, err := h.GetCurrentUserProfilePic(m.From.ID)
	if err != nil {
		return false, err
	}

	markedPic, err := watermarker.ApplyWatermark(ppic, pidormarkPicture)
	if err != nil {
		return false, err
	}

	err = h.PostReplyWithImage(m.Chat.ID, m.MessageID, markedPic, "pidormark.png", "Пидора ответ.")

	if err != nil {
		return false, err
	}
	return true, nil
}
