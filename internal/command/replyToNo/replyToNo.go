package replyToNo

import (
	"bytes"
	"context"
	_ "embed"
	"io"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var noRegex = regexp.MustCompile(`^(Net|net|Нет|нет)(.|\?|!)?$`)

const pidorText = "Пидора ответ."

//go:embed pidormark.png
var pidormarkPicture []byte

type ReplyToNo struct {
	messenger   messenger
	watermarker watermarker
	queue       queue
}

type watermarker interface {
	Apply(bakground io.Reader, foreground io.Reader, writer io.Writer) error
}

type messenger interface {
	ReplyWithImg(chatId int64, replyToMessageId int, img io.Reader, imgName string, caption string) (int, error)
	GetUserCurrentProfilePic(userId int64, w io.Writer) error
}

type queue interface {
	EnqueueDeleteMessage(chatId int64, messageId int) error
}

func New(messenger messenger, watermarker watermarker, queue queue) *ReplyToNo {
	return &ReplyToNo{
		messenger:   messenger,
		watermarker: watermarker,
		queue:       queue,
	}
}

func (h *ReplyToNo) Handle(ctx context.Context, m *tgbotapi.Message) error {
	ppicBuf := bytes.NewBuffer([]byte{})
	err := h.messenger.GetUserCurrentProfilePic(m.From.ID, ppicBuf)
	if err != nil {
		return err
	}

	resBuf := bytes.NewBuffer([]byte{})
	err = h.watermarker.Apply(ppicBuf, bytes.NewReader(pidormarkPicture), resBuf)
	if err != nil {
		return err
	}

	_, err = h.messenger.ReplyWithImg(m.Chat.ID, m.MessageID, resBuf, "pidormark.png", pidorText)

	if err != nil {
		return err
	}

	return nil
}

func (h *ReplyToNo) ShouldRun(m *tgbotapi.Message) bool {
	return noRegex.MatchString(m.Text)
}
