package tgcommand

import (
	"context"
	"fmt"

	"github.com/kapitan123/telegrofler/service/bot/domain/message"
)

type RecordReaction struct {
	storage postStorage
}

func NewRecordReaction(storage postStorage) *RecordReaction {
	return &RecordReaction{
		storage: storage,
	}
}

func (h *RecordReaction) Handle(ctx context.Context, m message.Message) error {
	existingPost, found, err := h.storage.GetPostById(ctx, m.MediaId())

	if err != nil {
		return err
	}

	// should never happen
	if !found {
		return fmt.Errorf("original post was not saved, reaction will not be recorded")
	}

	existingPost.AddReaction(m.AsReaction())

	h.storage.UpsertPost(ctx, existingPost)

	return nil
}

func (h *RecordReaction) ShouldRun(m message.Message) bool {
	return m.ReplytoMessage.HasMedia()
}
