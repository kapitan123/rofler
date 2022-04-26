package tghandler

import (
	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/data/post"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PostTopRoflerCommand struct {
	*bot.Bot
	*post.PostsStore
}

func NewPostTopRoflerCommand(b *bot.Bot, ps *post.PostsStore) *PostTopRoflerCommand {
	return &PostTopRoflerCommand{
		Bot:         b,
		PostsStore:  ps,
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
