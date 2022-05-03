package tgaction

import (
	"github.com/kapitan123/telegrofler/data/firestore"
	"github.com/kapitan123/telegrofler/internal/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type RecordBotPostReaction struct {
	*bot.Bot
	*firestore.PostsStore
}

func NewRecordBotPostReaction(b *bot.Bot, ps *firestore.PostsStore) *RecordBotPostReaction {
	return &RecordBotPostReaction{
		Bot:        b,
		PostsStore: ps,
	}
}

func (h *RecordBotPostReaction) Handle(m *tgbotapi.Message) (bool, error) {
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

	exPost, found, err := h.GetById(mediaRepy.VideoId)

	// in this case we don't record reaction as all bot posts should be saved already
	if !found {
		return isHandeled, nil
	}

	if err != nil {
		return isHandeled, err
	}

	exPost.AddReaction(details.Sender, details.Text, details.MessageId)
	h.Upsert(exPost)
	return isHandeled, nil
}
