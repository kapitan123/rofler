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
	GetPostById(ctx context.Context, videoId string) (storage.Post, bool, error)
	UpsertPost(ctx context.Context, p storage.Post) error
}

func New(messenger messenger, storage postStorage) *RecordBotPostReaction {
	return &RecordBotPostReaction{
		messenger: messenger,
		storage:   storage,
	}
}

func (h *RecordBotPostReaction) Handle(ctx context.Context, m *tgbotapi.Message) error {
	// AK TODO parsing the message twice rename package to parser
	reply := extractVideoRepostReaction(m)

	exPost, found, err := h.storage.GetPostById(ctx, reply.VideoId)

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

type replyToMediaPost struct {
	VideoId   string
	MessageId int // RepllyToMessage.ID not the update.Message.ID
	Sender    string
	Text      string
}

func extractVideoRepostReaction(upd *tgbotapi.Message) replyToMediaPost {
	rtm := upd.ReplyToMessage
	sender := upd.From.UserName

	reply := replyToMediaPost{
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

	// if the user reference his own post it is not a reaction.
	return sender != poster
}
