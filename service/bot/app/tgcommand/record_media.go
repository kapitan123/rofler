package tgcommand

import (
	"context"

	"github.com/kapitan123/telegrofler/common/logs"
	"github.com/kapitan123/telegrofler/service/bot/domain"
	"github.com/kapitan123/telegrofler/service/bot/domain/message"
)

type RecordMedia struct {
	storage postStorage
}

func NewRecordMediaPost(storage postStorage) *RecordMedia {
	return &RecordMedia{
		storage: storage,
	}
}

func (h *RecordMedia) Handle(ctx context.Context, m message.Message) (err error) {
	defer func() {
		logs.LogExecutionResult("RecordMedia ", m, err)
	}()

	newPost := domain.NewPost(m.MediaId(), m.From(), m.ChatId())
	err = h.storage.UpsertPost(ctx, newPost)

	if err != nil {
		return err
	}

	return nil
}

func (h *RecordMedia) ShouldRun(m message.Message) bool {
	return m.HasMedia() && !m.IsBotPost()
}
