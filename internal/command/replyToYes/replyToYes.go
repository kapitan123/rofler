package replyToYes

import (
	"bytes"
	"context"
	_ "embed"
	"io"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var yesRegex = regexp.MustCompile(`^(da|Да|да)(.|\?|!)?$`)

//go:embed kirkorov.png
var yesPicture []byte

type ReplyToYes struct {
	messenger messenger
}

type messenger interface {
	ReplyWithImg(chatId int64, replyToMessageId int, img io.Reader, imgName string, caption string) error
}

func New(messenger messenger) *ReplyToYes {
	return &ReplyToYes{
		messenger: messenger,
	}
}

func (h *ReplyToYes) Handle(ctx context.Context, m *tgbotapi.Message) error {
	err := h.messenger.ReplyWithImg(m.Chat.ID, m.MessageID, bytes.NewReader(yesPicture), "kirkorov.png", "")

	return err
}

func (h *ReplyToYes) ShouldRun(m *tgbotapi.Message) bool {
	return yesRegex.MatchString(m.Text)
}
