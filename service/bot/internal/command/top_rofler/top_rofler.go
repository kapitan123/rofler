package toprofler

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/service/bot/internal/messenger/format"
	"github.com/kapitan123/telegrofler/service/bot/internal/storage"
)

const commandName = "toprofler"

type TopRofler struct {
	messenger messenger
	storage   postStorage
}

type messenger interface {
	SendText(chatId int64, text string) (int, error)
}

type postStorage interface {
	GetChatPosts(ctx context.Context, chatId int64) ([]storage.Post, error)
}

func New(messenger messenger, storage postStorage) *TopRofler {
	return &TopRofler{
		messenger: messenger,
		storage:   storage,
	}
}

func (h *TopRofler) Handle(ctx context.Context, message *tgbotapi.Message) error {
	posts, err := h.storage.GetChatPosts(ctx, message.Chat.ID)
	if err != nil {
		return err
	}

	if len(posts) == 0 {
		return nil
	}

	roflerScores := countScores(posts)

	listMeassge := format.AsDescendingList(roflerScores, "🤡 <b>%s</b>: %d")

	_, err = h.messenger.SendText(message.Chat.ID, listMeassge)
	if err != nil {
		return err
	}
	return nil
}

func countScores(posts []storage.Post) map[string]int {
	roflerScores := map[storage.UserRef]int{}
	for _, p := range posts {
		roflerScores[p.UserRef] += len(p.Reactions)
	}

	names := map[string]int{}
	for k, v := range roflerScores {
		names[k.DisplayName] = v
	}
	return names
}

func (h *TopRofler) ShouldRun(message *tgbotapi.Message) bool {
	return message.IsCommand() && message.Command() == commandName
}
