package recordBotPostReaction

import (
	"context"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/storage"
	log "github.com/sirupsen/logrus"
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
	Upsert(ctx context.Context, p storage.Post) error
}

func New(messenger messenger, storage postStorage) *RecordBotPostReaction {
	return &RecordBotPostReaction{
		messenger: messenger,
		storage:   storage,
	}
}

func (h *RecordBotPostReaction) Handle(ctx context.Context, m *tgbotapi.Message) error {
	// AK TODO parsing the message twice rename package to parser
	reply, err := tryExtractVideoRepostReaction(m)

	if err != nil {
		return err
	}

	details := reply.Details
	exPost, found, err := h.storage.GetById(ctx, reply.VideoId)

	// in this case we don't record reaction as all bot posts should be saved already
	if !found {
		return nil
	}

	if err != nil {
		return err
	}

	exPost.AddReaction(details.Sender, details.Text, details.MessageId)
	h.storage.Upsert(ctx, exPost)
	return nil
}

func (h *RecordBotPostReaction) ShouldRun(m *tgbotapi.Message) bool {
	mediaReply, err := tryExtractVideoRepostReaction(m)

	if err != nil {
		log.WithError(err).Error("Reaction extraction has failed")
		return false
	}

	if mediaReply.Details.Sender == "" {
		return false
	}

	return true
}

func tryExtractVideoRepostReaction(upd *tgbotapi.Message) (bot.ReplyToMediaPost, error) {
	vr := bot.ReplyToMediaPost{}
	rtm := upd.ReplyToMessage

	if rtm == nil || rtm.From.UserName != "TelegroflBot" || rtm.Video == nil {
		return vr, nil
	}

	r := regexp.MustCompile(posterMaker)
	poster := r.FindStringSubmatch(rtm.Caption)[1]

	sender := upd.From.UserName
	// if the user reference his own post it is not a reaction. Can be moved outside the scope
	if sender == poster {
		return vr, nil
	}

	reply := bot.ReplyToMediaPost{
		VideoId: rtm.Video.FileName,
		Details: bot.Details{
			Sender:    sender,
			Text:      upd.Text,
			MessageId: upd.MessageID,
		},
	}

	return reply, nil
}
