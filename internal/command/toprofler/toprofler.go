package toprofler

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/storage"
)

const commandName = "top"

type TopRofler struct {
	messenger messenger
	storage   postStorage
}

type messenger interface {
	SendMessage(chatId int64, text string) error
}

type postStorage interface {
	GetAll(ctx context.Context) ([]storage.Post, error)
}

func New(messenger messenger, storage postStorage) *TopRofler {
	return &TopRofler{
		messenger: messenger,
		storage:   storage,
	}
}

func (h *TopRofler) Handle(ctx context.Context, message *tgbotapi.Message) error {
	posts, err := h.storage.GetAll(ctx)
	if err != nil {
		return err
	}

	roflerScores := countScores(posts)

	listMeassge := formatListMessage(roflerScores)

	err = h.messenger.SendMessage(message.Chat.ID, listMeassge)
	if err != nil {
		return err
	}
	return nil
}

func formatListMessage(roflerScores map[string]int) string {
	listMeassge := ""
	for name, score := range roflerScores {
		listMeassge += formatTopRofler(name, score)
	}
	return listMeassge
}

func countScores(posts []storage.Post) map[string]int {
	roflerScores := map[string]int{}
	for _, p := range posts {
		roflerScores[p.RoflerUserName] += len(p.Reactions)
	}
	return roflerScores
}

func (h *TopRofler) ShouldRun(message *tgbotapi.Message) bool {
	return message.IsCommand() && message.Command() == commandName
}

func formatTopRofler(username string, score int) string {
	return fmt.Sprintf("\U0001F451 <b>@%s</b> <b>Likes:</b> %d", username, score)
}
