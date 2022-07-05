package recordBotPostReaction

import (
	"context"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/storage"
)

const posterMaker = `🔥@(.*?)🔥`

type RecordBotPostReaction struct {
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

type (
	ReplyToMediaPost struct {
		VideoId   string
		MessageId int // RepllyToMessage.ID not the update.Message.ID
		Sender    string
		Text      string
	}
)

func New(messenger messenger, storage postStorage) *RecordBotPostReaction {
	return &RecordBotPostReaction{
		messenger: messenger,
		storage:   storage,
	}
}

func (h *RecordBotPostReaction) Handle(ctx context.Context, m *tgbotapi.Message) error {
	// AK TODO parsing the message twice rename package to parser
	reply := extractVideoRepostReaction(m)

	exPost, found, err := h.storage.GetById(ctx, reply.VideoId)

	// in this case we don't record reaction as all bot posts should be saved already
	if !found {
		return nil
	}

	if err != nil {
		return err
	}

	exPost.AddReaction(reply.Sender, reply.Text, reply.MessageId)
	h.storage.UpsertPost(ctx, exPost)
	return nil
}

func (h *RecordBotPostReaction) ShouldRun(m *tgbotapi.Message) bool {
	return containsVideoRepostReaction(m)
}

func extractVideoRepostReaction(upd *tgbotapi.Message) ReplyToMediaPost {
	rtm := upd.ReplyToMessage
	sender := upd.From.UserName

	reply := ReplyToMediaPost{
		VideoId:   rtm.Video.FileName,
		Sender:    sender,
		Text:      upd.Text,
		MessageId: upd.MessageID,
	}

	return reply
}

func containsVideoRepostReaction(upd *tgbotapi.Message) bool {
	rtm := upd.ReplyToMessage

	if rtm == nil || rtm.From.UserName != "TelegroflBot" || rtm.Video == nil {
		return false
	}

	r := regexp.MustCompile(posterMaker)
	poster := r.FindStringSubmatch(rtm.Caption)[1]

	sender := upd.From.UserName
	// if the user reference his own post it is not a reaction. Can be moved outside the scope
	if sender == poster {
		return false
	}

	return true
}
