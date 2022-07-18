package recordBotPostReaction

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/storage"
)

type RecordBotPostReaction struct {
	messenger messenger
	storage   postStorage
}

type (
	messenger interface {
		ReplyWithText(chatId int64, messageId int, text string) error
	}

	postStorage interface {
		GetPostById(ctx context.Context, videoId string) (storage.Post, bool, error)
		UpsertPost(ctx context.Context, p storage.Post) error
	}
)

func New(messenger messenger, storage postStorage) *RecordBotPostReaction {
	return &RecordBotPostReaction{
		messenger: messenger,
		storage:   storage,
	}
}

func (h *RecordBotPostReaction) Handle(ctx context.Context, m *tgbotapi.Message) error {
	reply := extractVideoRepostReaction(m)

	exPost, found, err := h.storage.GetPostById(ctx, reply.VideoId)

	if err != nil {
		return err
	}

	// in this case we don't record reaction as all bot posts should be saved already
	if !found {
		return nil
	}

	exPost.AddReaction(reply.ReactorRef, reply.Text, reply.ToMessageId)
	h.storage.UpsertPost(ctx, exPost)
	return nil
}

func (h *RecordBotPostReaction) ShouldRun(m *tgbotapi.Message) bool {
	return containsVideoRepostReaction(m)
}

type replyToMediaPost struct {
	VideoId     string
	ToMessageId int // RepllyToMessage.ID not the update.Message.ID
	ReactorRef  storage.UserRef
	Text        string
}

func extractVideoRepostReaction(upd *tgbotapi.Message) replyToMediaPost {
	rtm := upd.ReplyToMessage
	reply := replyToMediaPost{
		VideoId: rtm.Video.FileName,
		ReactorRef: storage.UserRef{
			Id:          rtm.From.ID,
			DisplayName: fmt.Sprintf("%s %s", upd.From.FirstName, upd.From.LastName),
		},
		Text:        upd.Text,
		ToMessageId: rtm.MessageID,
	}

	return reply
}

func containsVideoRepostReaction(upd *tgbotapi.Message) bool {
	rtm := upd.ReplyToMessage

	if rtm == nil || rtm.From.UserName != "TelegroflBot" || rtm.Video == nil {
		return false
	}

	// monitored bot posts has exactly one mention
	if len(rtm.Entities) == 0 {
		return false
	}

	// if the user reference his own post it is not a reaction.
	return upd.From.ID != rtm.Entities[0].User.ID
}
