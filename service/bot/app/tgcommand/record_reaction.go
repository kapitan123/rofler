package tgcommand

import (
	"context"

	"github.com/kapitan123/telegrofler/common/logs"
	"github.com/kapitan123/telegrofler/service/bot/domain"
)

type RecordReaction struct {
	storage postStorage
}

func NewRecordReaction(storage postStorage) *RecordReaction {
	return &RecordReaction{
		storage: storage,
	}
}

func (h *RecordReaction) Handle(ctx context.Context, m domain.Message) (err error) {
	defer func() {
		logs.LogExecutionResult("RecordReaction ", m, err)
	}()

	// existingPost, found, err := h.storage.GetPostById(ctx, m.MediaId())

	// if err != nil {
	// 	return err
	// }

	// // should never happen
	// if !found {
	// 	return fmt.Errorf("original post was not saved, reaction will not be recorded")
	// }

	// existingPost.AddReaction(m.AsReaction())

	// h.storage.UpsertPost(ctx, existingPost)

	return nil
}

func (h *RecordReaction) ShouldRun(m domain.Message) bool {
	return m.ReplytoMessage.HasMedia()
}
