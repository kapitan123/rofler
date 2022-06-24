package tgaction

import (
	"context"

	"github.com/kapitan123/telegrofler/data/firestore/posts"
	"github.com/kapitan123/telegrofler/internal/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PostTopRoflerCommand struct {
	*bot.Bot
	postsStorage *posts.PostsStorage
}

func NewPostTopRoflerCommand(b *bot.Bot, ps *posts.PostsStorage) *PostTopRoflerCommand {
	return &PostTopRoflerCommand{
		Bot:          b,
		postsStorage: ps,
	}
}

func (h *PostTopRoflerCommand) Handle(m *tgbotapi.Message, ctx context.Context) (bool, error) {
	tr, roflCount, err := h.postsStorage.GetTopRoflerFromPosts(ctx)
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
