package toppidor

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/service/bot/internal/messenger/format"
	"github.com/kapitan123/telegrofler/service/bot/internal/storage"
)

const commandName = "toppidor"

type TopPidor struct {
	messenger messenger
	storage   postStorage
}

type (
	messenger interface {
		SendText(chatId int64, text string) (int, error)
	}

	postStorage interface {
		GetChatPidors(ctx context.Context, chatId int64) ([]storage.Pidor, error)
	}
)

func New(messenger messenger, storage postStorage) *TopPidor {
	return &TopPidor{
		messenger: messenger,
		storage:   storage,
	}
}

func (h *TopPidor) Handle(ctx context.Context, message *tgbotapi.Message) error {
	pidors, err := h.storage.GetChatPidors(ctx, message.Chat.ID)
	if err != nil {
		return err
	}

	pidorScores := countScores(pidors)

	listMeassge := format.AsDescendingList(pidorScores, "🐓 <b>%s:</b> %d ⚣⚣")

	_, err = h.messenger.SendText(message.Chat.ID, listMeassge)
	if err != nil {
		return err
	}
	return nil
}

func countScores(pidors []storage.Pidor) map[string]int {
	pidorScores := map[storage.UserRef]int{}
	for _, p := range pidors {
		pidorScores[p.UserRef] += 1
	}

	names := map[string]int{}
	for k, v := range pidorScores {
		names[k.DisplayName] = v
	}
	return names
}

func (h *TopPidor) ShouldRun(message *tgbotapi.Message) bool {
	return message.IsCommand() && message.Command() == commandName
}
