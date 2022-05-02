package tgaction

import (
	_ "embed"
	"regexp"

	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/pkg/watermarker"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var noRegex = regexp.MustCompile(`^(Net|net|Нет|нет)(.|\?|!)?$`)

const pidorText = "Пидора ответ."

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
	found := noRegex.MatchString(m.Text)

	if !found {
		return false, nil
	}

	ppic, _ := h.GetCurrentUserProfilePic(m.From.ID)

	if ppic == nil {
		h.PostReplyWithText(m.Chat.ID, m.MessageID, pidorText)
		return true, nil
	} else {
		markedPic, err := watermarker.ApplyWatermark(ppic, pidormarkPicture)
		if err != nil {
			return true, err
		}

		err = h.PostReplyWithImage(m.Chat.ID, m.MessageID, markedPic, "pidormark.png", pidorText)

		if err != nil {
			return true, err
		}
		return true, nil
	}
}
