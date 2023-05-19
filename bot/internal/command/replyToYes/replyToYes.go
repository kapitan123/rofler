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
	queue     queue
}

type queue interface {
	EnqueueDeleteMessage(chatId int64, messageId int) error
}

type messenger interface {
	ReplyWithImg(chatId int64, replyToMessageId int, img io.Reader, imgName string, caption string) (int, error)
}

func New(messenger messenger, queue queue) *ReplyToYes {
	return &ReplyToYes{
		messenger: messenger,
		queue:     queue,
	}
}

func (h *ReplyToYes) Handle(ctx context.Context, m *tgbotapi.Message) error {
	chatId := m.Chat.ID
	newMessageId, err := h.messenger.ReplyWithImg(chatId, m.MessageID, bytes.NewReader(yesPicture), "kirkorov.png", "")

	if err != nil {
		return err
	}

	return h.queue.EnqueueDeleteMessage(chatId, newMessageId)
}

func (h *ReplyToYes) ShouldRun(m *tgbotapi.Message) bool {
	return yesRegex.MatchString(m.Text)
}
