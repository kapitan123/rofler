package tgaction

import (
	"context"
	"time"

	"github.com/kapitan123/telegrofler/data/firestore/posts"
	"github.com/kapitan123/telegrofler/internal/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type RecordReactionToUserMediaPost struct {
	*bot.Bot
	postsStorage *posts.PostsStorage
}

func NewRecordReactionToUserMediaPost(b *bot.Bot, ps *posts.PostsStorage) *RecordReactionToUserMediaPost {
	return &RecordReactionToUserMediaPost{
		Bot:          b,
		postsStorage: ps,
	}
}

func (h *RecordReactionToUserMediaPost) Handle(m *tgbotapi.Message, ctx context.Context) (bool, error) {
	rtm := m.ReplyToMessage

	if rtm == nil || rtm.Video == nil {
		return false, nil
	}

	isHandeled := true

	mediaRepy, err := bot.ExtractUserMediaReaction(m)
	details := mediaRepy.Details
	if err != nil {
		return !isHandeled, err
	}

	// AK TODO should actually return nil
	if details.Sender == "" {
		return !isHandeled, nil
	}

	log.Infof("Reaction was found for %s sent by %s", mediaRepy.VideoId, details.Sender)

	exPost, found, err := h.postsStorage.GetById(ctx, mediaRepy.VideoId)

	if err != nil {
		return isHandeled, err
	}

	if !found {
		reactions := make([]posts.Reaction, 0)
		exPost = posts.Post{
			VideoId:        mediaRepy.VideoId,
			Source:         "misc",
			RoflerUserName: rtm.From.UserName,
			Url:            "",
			Reactions:      reactions,
			PostedOn:       time.Now(),
		}
	}

	exPost.AddReaction(details.Sender, details.Text, details.MessageId)
	h.postsStorage.Upsert(ctx, exPost)

	return isHandeled, nil
}
