package tgaction

import (
	"context"
	"cloud.google.com/go/firestore"

	"github.com/kapitan123/telegrofler/data/firestore/posts"
	"github.com/kapitan123/telegrofler/internal/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type RecordBotPostReaction struct {
	*bot.Bot
	FsClient *firestore.Client
}

func NewRecordBotPostReaction(b *bot.Bot, fs *firestore.Client) *RecordBotPostReaction {
	return &RecordBotPostReaction{
		Bot:        b,
		FsClient: fs,
	}
}

func (h *RecordBotPostReaction) Handle(m *tgbotapi.Message, ctx context.Context) (bool, error) {
	isHandeled := true
	mediaRepy, err := bot.TryExtractVideoRepostReaction(m)
	details := mediaRepy.Details

	if err != nil {
		return !isHandeled, err
	}

	if details.Sender == "" {
		return !isHandeled, nil
	}

	log.Infof("Reaction was found for %s sent by %s", mediaRepy.VideoId, details.Sender)

	exPost, found, err := posts.GetById(ctx, h.FsClient, mediaRepy.VideoId)

	// in this case we don't record reaction as all bot posts should be saved already
	if !found {
		return isHandeled, nil
	}

	if err != nil {
		return isHandeled, err
	}

	exPost.AddReaction(details.Sender, details.Text, details.MessageId)
	posts.Upsert(ctx, h.FsClient, exPost)
	return isHandeled, nil
}
