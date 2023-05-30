package tgcommand

import (
	"context"

	"github.com/kapitan123/telegrofler/service/bot/domain"
	"github.com/kapitan123/telegrofler/service/bot/domain/format"
	"github.com/kapitan123/telegrofler/service/bot/domain/message"
)

const commandName = "toprofler"

type TopRofler struct {
	messenger messenger
	storage   postStorage
}

func NewTopRofler(messenger messenger, storage postStorage) *TopRofler {
	return &TopRofler{
		messenger: messenger,
		storage:   storage,
	}
}

func (h *TopRofler) Handle(ctx context.Context, message message.Message) error {
	posts, err := h.storage.GetChatPosts(ctx, message.ChatId())
	if err != nil {
		return err
	}

	if len(posts) == 0 {
		return nil
	}

	// AK TODO also should be gone to domain
	roflerScores := countScores(posts)

	listMeassge := format.AsDescendingList(roflerScores, "🤡 <b>%s</b>: %d")

	_, err = h.messenger.SendText(message.ChatId(), listMeassge)
	if err != nil {
		return err
	}
	return nil
}

func countScores(posts []domain.Post) map[string]int {
	roflerScores := map[domain.UserRef]int{}
	for _, p := range posts {
		roflerScores[p.Poster] += len(p.Reactions)
	}

	names := map[string]int{}
	for k, v := range roflerScores {
		names[k.DisplayName] = v
	}
	return names
}

func (h *TopRofler) ShouldRun(message message.Message) bool {
	return message.IsCommand(commandName)
}
