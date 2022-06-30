package replyToNo

import (
	"context"
	_ "embed"
	"regexp"

	"github.com/kapitan123/telegrofler/internal/watermark"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var noRegex = regexp.MustCompile(`^(Net|net|Нет|нет)(.|\?|!)?$`)

const pidorText = "Пидора ответ."

//go:embed pidormark.png
var pidormarkPicture []byte

type ReplyToNo struct {
	messenger messenger
}

type messenger interface {
	ReplyWithImg(chatId int64, replyToMessageId int, img []byte, imgName string, caption string) error
	GetUserCurrentProfilePic(userId int64) ([]byte, error)
}

func New(messenger messenger) *ReplyToNo {
	return &ReplyToNo{
		messenger: messenger,
	}
}

func (h *ReplyToNo) Handle(ctx context.Context, m *tgbotapi.Message) error {
	ppic, _ := h.messenger.GetUserCurrentProfilePic(m.From.ID)

	markedPic, err := watermark.Apply(ppic, pidormarkPicture)
	if err != nil {
		return err
	}

	err = h.messenger.ReplyWithImg(m.Chat.ID, m.MessageID, markedPic, "pidormark.png", pidorText)

	if err != nil {
		return err
	}

	return nil
}

func (h *ReplyToNo) ShouldRun(m *tgbotapi.Message) bool {
	return noRegex.MatchString(m.Text)
}
