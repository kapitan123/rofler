package recordReaction

import (
	"context"

	"github.com/kapitan123/telegrofler/service/bot/domain"
	"github.com/kapitan123/telegrofler/service/bot/domain/message"
)

// AK TODO rewrite this lil boi
type RecordReaction struct {
	messenger messenger
	storage   postStorage
}

type (
	messenger interface {
		ReplyWithText(chatId int64, messageId int, text string) (int, error)
	}

	postStorage interface {
		GetPostById(ctx context.Context, videoId string) (domain.Post, bool, error)
		UpsertPost(ctx context.Context, p domain.Post) error
	}
)

func New(messenger messenger, storage postStorage) *RecordReaction {
	return &RecordReaction{
		messenger: messenger,
		storage:   storage,
	}
}

func (h *RecordReaction) Handle(ctx context.Context, m message.Message) error {
	post, found, err := h.storage.GetPostById(ctx, m.MediaId())

	if err != nil {
		return err
	}

	if !found {
		post = domain.NewPost(m.ReplyToMessage.From(), m.ReplyToMessage.ChatId())
	}

	post.AddReaction(m.AsReaction())
	h.storage.UpsertPost(ctx, post)

	return nil
}

func (h *RecordReaction) ShouldRun(m message.Message) bool {
	return m.IsReplyToMedia() && !m.IsSelfReply() && !m.IsBotPost()
}
