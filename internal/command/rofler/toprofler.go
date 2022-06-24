package rofler

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/storage"
)

type TopRofler struct {
	messenger messenger
	storage   postStorage
}

type messenger interface {
	SendMessage(ctx context.Context, chatId int64, text string) error
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

	roflerScores := map[string]int{}
	for _, p := range posts {
		roflerScores[p.RoflerUserName] += len(p.Reactions)
	}

	maxUserName, max := "", 0
	for username, score := range roflerScores {
		if max < score {
			max = score
			maxUserName = username
		}
	}
	if maxUserName == "" {
		return nil
	}
	err = h.messenger.SendMessage(ctx, message.Chat.ID, formatTopRofler(maxUserName, max))
	if err != nil {
		return err
	}
	return nil
}

func (h *TopRofler) ShouldRun(message *tgbotapi.Message) bool {
	return message.IsCommand() && message.Command() == "top"
}

func formatTopRofler(username string, score int) string {
	return fmt.Sprintf("\U0001F451 <b>@%s</b>\n<b>Likes:</b> %d", username, score)
}
