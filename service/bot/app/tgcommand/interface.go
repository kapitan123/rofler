package tgcommand

import (
	"context"
	"io"

	"github.com/kapitan123/telegrofler/service/bot/domain"
)

type messenger interface {
	Delete(chatId domain.ChatId, messageId int) error
	SendText(chatId domain.ChatId, text string) (int, error)
}

type postStorage interface {
	GetPostById(ctx context.Context, videoId string) (domain.Post, bool, error)
	UpsertPost(ctx context.Context, p domain.Post) error
	GetChatPosts(ctx context.Context, chatId domain.ChatId) ([]domain.Post, error)
}

type filesBucket interface {
	Read(ctx context.Context, addr string) (io.Reader, error)
}
