package tgaction

import (
	"github.com/kapitan123/telegrofler/data/firestore/posts"
	"github.com/kapitan123/telegrofler/internal/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PostTopRoflerCommand struct {
	*bot.Bot
	*firestore.PostsStore
}

func NewPostTopRoflerCommand(b *bot.Bot, ps *firestore.PostsStore) *PostTopRoflerCommand {
	return &PostTopRoflerCommand{
		Bot:        b,
		PostsStore: ps,
	}
}

func (h *PostTopRoflerCommand) Handle(m *tgbotapi.Message) (bool, error) {
	tr, roflCount, err := h.GetTopRoflerFromPosts()
	if err != nil {
		return false, err
	}

	err = h.PostTopRofler(m.Chat.ID, tr, roflCount)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (h *PostTopRoflerCommand) GetCommandText() string {
	return "top"
}
