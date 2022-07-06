package toppidor

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/storage"
)

const commandName = "toppidor"

type TopPidor struct {
	messenger messenger
	storage   postStorage
	formatter formatter
}

type (
	formatter interface {
		FormatAsDescendingList(map[string]int, string) string
	}

	messenger interface {
		SendText(chatId int64, text string) error
	}

	postStorage interface {
		GetAllPidors(ctx context.Context) ([]storage.Pidor, error)
	}
)

func New(messenger messenger, storage postStorage, formatter formatter) *TopPidor {
	return &TopPidor{
		messenger: messenger,
		storage:   storage,
		formatter: formatter,
	}
}

func (h *TopPidor) Handle(ctx context.Context, message *tgbotapi.Message) error {
	pidors, err := h.storage.GetAllPidors(ctx)
	if err != nil {
		return err
	}

	pidorScores := countScores(pidors)

	listMeassge := h.formatter.FormatAsDescendingList(pidorScores, "🐓 <b>%s</b> <b>was pidor:</b> %d times ⚣⚣")

	err = h.messenger.SendText(message.Chat.ID, listMeassge)
	if err != nil {
		return err
	}
	return nil
}

func countScores(posts []storage.Pidor) map[string]int {
	roflerScores := map[string]int{}
	for _, p := range posts {
		roflerScores[p.UserName] += 1
	}
	return roflerScores
}

func (h *TopPidor) ShouldRun(message *tgbotapi.Message) bool {
	return message.IsCommand() && message.Command() == commandName
}
