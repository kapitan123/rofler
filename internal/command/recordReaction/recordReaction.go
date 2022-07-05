package recordReaction

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/storage"
	log "github.com/sirupsen/logrus"
)

type RecordReaction struct {
	messenger messenger
	storage   postStorage
}

type messenger interface {
	ReplyWithText(chatId int64, messageId int, text string) error
}

type postStorage interface {
	GetById(ctx context.Context, videoId string) (storage.Post, bool, error)
	UpsertPost(ctx context.Context, p storage.Post) error
}

func New(messenger messenger, storage postStorage) *RecordReaction {
	return &RecordReaction{
		messenger: messenger,
		storage:   storage,
	}
}

func (h *RecordReaction) Handle(ctx context.Context, m *tgbotapi.Message) error {
	mediaReply := extractUserMediaReaction(m)

	log.Infof("Reaction was found for %s sent by %s", mediaReply.VideoId, mediaReply.Details.Sender)

	details := mediaReply.Details
	exPost, found, err := h.storage.GetById(ctx, mediaReply.VideoId)

	if err != nil {
		return err
	}

	rtm := m.ReplyToMessage

	if !found {
		reactions := make([]storage.Reaction, 0)
		exPost = storage.Post{
			VideoId:        mediaReply.VideoId,
			Source:         "misc",
			RoflerUserName: rtm.From.UserName,
			Url:            "",
			Reactions:      reactions,
			PostedOn:       time.Now(),
		}
	}

	exPost.AddReaction(details.Sender, details.Text, details.MessageId)
	h.storage.UpsertPost(ctx, exPost)

	return nil
}

func (h *RecordReaction) ShouldRun(m *tgbotapi.Message) bool {
	rtm := m.ReplyToMessage

	if rtm == nil || rtm.Video == nil || m.From.UserName == "" {
		return false
	}

	return true
}

func extractUserMediaReaction(upd *tgbotapi.Message) ReplyToMediaPost {
	rtm := upd.ReplyToMessage
	vr := ReplyToMediaPost{
		VideoId: rtm.Video.FileID,
		Details: Details{
			Sender:    upd.From.UserName,
			MessageId: rtm.MessageID,
			Text:      upd.Text,
		},
	}

	return vr
}

type (
	ReplyToMediaPost struct {
		VideoId string
		Details Details
	}

	Details struct {
		MessageId int // RepllyToMessage.ID not the update.Message.ID
		Sender    string
		Text      string
	}
)
