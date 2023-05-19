package choosePidor

import (
	"bytes"
	"context"
	"crypto/rand"
	_ "embed"
	"fmt"
	"io"
	"math/big"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/bot/internal/messenger/format"
	"github.com/kapitan123/telegrofler/bot/internal/storage"
	log "github.com/sirupsen/logrus"
)

//go:embed pidormark.png
var pidormarkPicture []byte

//go:embed tinfoil.jpg
var tinfoilPicture []byte

const commandName = "choosepidor"

type ChoosePidor struct {
	messenger   messenger
	storage     pidorStorage
	watermarker watermarker
	systemclock systemclock
	queue       queue
}

type (
	watermarker interface {
		Apply(bakground io.Reader, foreground io.Reader, w io.Writer) error
	}

	messenger interface {
		SendText(chatId int64, text string) (int, error)
		SendImg(chatId int64, img io.Reader, imgName string, caption string) (int, error)
		GetChatAdmins(chatId int64) ([]tgbotapi.ChatMember, error)
		GetUserCurrentProfilePic(userId int64, w io.Writer) error
	}

	pidorStorage interface {
		GetPidorForDate(ctx context.Context, chatid int64, date time.Time) (storage.Pidor, bool, error)
		CreatePidor(ctx context.Context, chatid int64, userRef storage.UserRef, date time.Time) error
	}

	systemclock interface {
		Now() time.Time
	}

	queue interface {
		EnqueueDeleteMessage(chatId int64, messageId int) error
	}
)

func New(messenger messenger, storage pidorStorage, watermarker watermarker, queue queue, systemclock systemclock) *ChoosePidor {
	return &ChoosePidor{
		messenger:   messenger,
		storage:     storage,
		watermarker: watermarker,
		systemclock: systemclock,
		queue:       queue,
	}
}

func (h *ChoosePidor) Handle(ctx context.Context, m *tgbotapi.Message) error {
	return h.ChoosePidor(ctx, m.MessageID, m.Chat.ID)
}

func (h *ChoosePidor) ChoosePidor(ctx context.Context, msgId int, chatId int64) error {
	currDate := h.systemclock.Now()
	pidor, found, err := h.storage.GetPidorForDate(ctx, chatId, currDate)

	if err != nil {
		return err
	}

	if msgId > 0 {
		err = h.queue.EnqueueDeleteMessage(chatId, msgId)

		if err != nil {
			return err
		}
	}

	if found {
		mention := format.AsUserMention(pidor.UserRef.Id, pidor.UserRef.DisplayName)
		newMsgId, err := h.messenger.SendText(chatId, fmt.Sprintf(mention+" is still sucking juicy cocks"))

		if err != nil {
			return err
		}

		return h.queue.EnqueueDeleteMessage(chatId, newMsgId)
	}

	admins, err := h.messenger.GetChatAdmins(chatId)

	if err != nil {
		return err
	}

	chosenOne := chooseRandom(admins).User

	fullName := chosenOne.FirstName + " " + chosenOne.LastName
	ur := storage.UserRef{Id: chosenOne.ID, DisplayName: fullName}

	err = h.storage.CreatePidor(ctx, chatId, ur, currDate)
	if err != nil {
		return err
	}

	ppicBuf := bytes.NewBuffer([]byte{})
	err = h.messenger.GetUserCurrentProfilePic(chosenOne.ID, ppicBuf)

	mention := format.AsUserMention(ur.Id, ur.DisplayName)

	if err != nil {
		log.WithError(err).Error("failed to generate user profile pic")
		_, err := h.messenger.SendImg(chatId, bytes.NewReader(tinfoilPicture), "tinfoil.png", "Скрытный пидор дня у нас "+mention)
		return err
	}

	resBuf := bytes.NewBuffer([]byte{})
	err = h.watermarker.Apply(ppicBuf, bytes.NewReader(pidormarkPicture), resBuf)

	if err != nil {
		return err
	}

	_, err = h.messenger.SendImg(chatId, resBuf, "pidor.png", "Pidor of the day is "+mention)
	return err
}

func chooseRandom(members []tgbotapi.ChatMember) tgbotapi.ChatMember {
	randIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(members)-1)))
	return members[randIndex.Int64()]
}

func (h *ChoosePidor) ShouldRun(message *tgbotapi.Message) bool {
	return message != nil && message.IsCommand() && message.Command() == commandName
}
