package toppidor

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/helpers/sortedmap"
	"github.com/kapitan123/telegrofler/internal/storage"
)

const commandName = "toppidor"

type TopPidor struct {
	messenger messenger
	storage   postStorage
}

type messenger interface {
	SendText(chatId int64, text string) error
}

type postStorage interface {
	GetAllPidors(ctx context.Context) ([]storage.Pidor, error)
}

func New(messenger messenger, storage postStorage) *TopPidor {
	return &TopPidor{
		messenger: messenger,
		storage:   storage,
	}
}

func (h *TopPidor) Handle(ctx context.Context, message *tgbotapi.Message) error {
	pidors, err := h.storage.GetAllPidors(ctx)
	if err != nil {
		return err
	}

	roflerScores := countScores(pidors)

	sortedRoflerScores := sortedmap.Sort(roflerScores)
	listMeassge := formatListMessage(sortedRoflerScores)

	err = h.messenger.SendText(message.Chat.ID, listMeassge)
	if err != nil {
		return err
	}
	return nil
}

// AK TODO shoud be extracted to a list helper
func formatListMessage(roflerScores sortedmap.PairList) string {
	listMeassge := ""
	for _, pair := range roflerScores {
		listMeassge += formatTopRofler(pair.Key, pair.Value)
	}
	return listMeassge
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

func formatTopRofler(username string, score int) string {
	return fmt.Sprintf("🐓 <b>%s</b> <b>was pidor:</b> %d times ⚣⚣", username, score)
}
