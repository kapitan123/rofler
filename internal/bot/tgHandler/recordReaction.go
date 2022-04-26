package tghandler

import (
	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/data/post"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type RecordReaction struct {
	*bot.Bot
	*post.PostsStore
}

func NewRecordReaction(b *bot.Bot, ps *post.PostsStore) *RecordReaction {
	return &RecordReaction{
		Bot:        b,
		PostsStore: ps,
	}
}

func (h *RecordReaction) Handle(m *tgbotapi.Message) (bool, error) {
	isHandeled := true
	reaction, err := h.TryExtractVideoRepostReaction(m)
	if err != nil {
		return !isHandeled, err
	}

	if reaction.Sender == "" {
		return isHandeled, nil
	}

	log.Infof("Reaction was found for %s sent by %s", reaction.VideoId, reaction.Sender)

	h.AddReactionToPost(reaction)

	return isHandeled, nil
}
