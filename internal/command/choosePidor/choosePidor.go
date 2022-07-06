package choosePidor

import (
	"context"
	_ "embed"
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/storage"
)

// AK TODO remove duplicate embedding
//go:embed pidormark.png
var pidormarkPicture []byte

const commandName = "choosePidor"

type ChoosePidor struct {
	messenger   messenger
	storage     pidorStorage
	watermarker watermarker
	systemclock systemclock
}

type (
	watermarker interface {
		Apply(bakground []byte, foreground []byte) ([]byte, error)
	}

	messenger interface {
		SendText(chatId int64, text string) error
		SendImg(chatId int64, img []byte, imgName string, caption string) error
		GetAdminUserNames(chatId int64) ([]string, error)
		GetUserCurrentProfilePic(userId int64) ([]byte, error)
	}

	pidorStorage interface {
		GetPidorForDate(ctx context.Context, chatid int64, date time.Time) (storage.Pidor, bool, error)
		CreatePidor(ctx context.Context, chatid int64, username string, date time.Time) error
	}

	systemclock interface {
		CurrentDate() time.Time
	}
)

func New(messenger messenger, storage pidorStorage, watermarker watermarker, systemclock systemclock) *ChoosePidor {
	return &ChoosePidor{
		messenger:   messenger,
		storage:     storage,
		watermarker: watermarker,
		systemclock: systemclock,
	}
}

func (h *ChoosePidor) Handle(ctx context.Context, m *tgbotapi.Message) error {
	currDate := h.systemclock.CurrentDate()
	pidor, found, err := h.storage.GetPidorForDate(ctx, m.Chat.ID, currDate)

	if err != nil {
		return err
	}

	if found {
		err = h.messenger.SendText(m.Chat.ID, pidor.UserName+" is still sucking juicy cocks")
		return err
	}

	names, err := h.messenger.GetAdminUserNames(m.Chat.ID)

	if err != nil {
		return err
	}

	todayPidorName := chooseRandom(names)

	err = h.storage.CreatePidor(ctx, m.Chat.ID, todayPidorName, currDate)
	if err != nil {
		return err
	}

	ppic, _ := h.messenger.GetUserCurrentProfilePic(m.From.ID)

	markedPic, err := h.watermarker.Apply(ppic, pidormarkPicture)
	if err != nil {
		return err
	}

	return h.messenger.SendImg(m.Chat.ID, markedPic, "pidor.png", "Pidor of the day is "+todayPidorName)
}

func chooseRandom(names []string) string {
	return names[rand.Intn(len(names)-1)]
}

func (h *ChoosePidor) ShouldRun(message *tgbotapi.Message) bool {
	return message.IsCommand() && message.Command() == commandName
}
