package tgaction

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/kapitan123/telegrofler/data/firestore/posts"
	"github.com/kapitan123/telegrofler/internal/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PostTopRoflerCommand struct {
	*bot.Bot
	FsClient *firestore.Client
}

func NewPostTopRoflerCommand(b *bot.Bot, fs *firestore.Client) *PostTopRoflerCommand {
	return &PostTopRoflerCommand{
		Bot:      b,
		FsClient: fs,
	}
}

func (h *PostTopRoflerCommand) Handle(m *tgbotapi.Message, ctx context.Context) (bool, error) {
	tr, roflCount, err := posts.GetTopRoflerFromPosts(ctx, h.FsClient)
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
