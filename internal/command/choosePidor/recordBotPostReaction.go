package recordBotPostReaction

import (
	"context"
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/storage"
)

// AK TODO remove duplicate embedding
//go:embed pidormark.png
var pidormarkPicture []byte

type ChoosePidor struct {
	messenger messenger
	storage   pidorStorage
}

type messenger interface {
	SendText(chatId int64, text string) error
	SendImg(chatId int64, img []byte, imgName string, caption string) error
	GetAdminUserNamess(chatId int64) ([]string, error)
}

type pidorStorage interface {
	GetForDate(ctx context.Context, date time.Time) (storage.Pidor, bool, error)
	CreatePidor(ctx context.Context, p storage.Pidor) error
}

func New(messenger messenger) *ChoosePidor {
	return &ChoosePidor{
		messenger: messenger,
	}
}

func (h *ChoosePidor) Handle(ctx context.Context, m *tgbotapi.Message) error {
	now := time.Now()
	pidor, found, err := h.storage.GetForDate(ctx, now)

	if err != nil {
		return err
	}

	if found {
		err = h.messenger.SendText(m.Chat.ID, "Pidor is still "+pidor.UserName)
		return nil
	}

	names, err := h.messenger.GetAdminUserNamess(m.Chat.ID)

	if err != nil {
		return err
	}

	todayPidorName := chooseRandom(names)

	err = h.storage.CreatePidor(ctx, storage.Pidor{ChosenOn: now, UserName: todayPidorName})
	if err != nil {
		return err
	}

	return h.messenger.SendImg(m.Chat.ID, pidormarkPicture, "pidor.png", "Pidor of the day is "+todayPidorName)
}

func chooseRandom(names []string) string {
	return names[rand.Intn(len(names)-1)]
}
