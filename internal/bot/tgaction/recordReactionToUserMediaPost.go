package tgaction

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/kapitan123/telegrofler/data/firestore/posts"
	"github.com/kapitan123/telegrofler/internal/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type RecordReactionToUserMediaPost struct {
	*bot.Bot
	FsClient *firestore.Client
}

func NewRecordReactionToUserMediaPost(b *bot.Bot, fs *firestore.Client) *RecordReactionToUserMediaPost {
	return &RecordReactionToUserMediaPost{
		Bot:      b,
		FsClient: fs,
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

	exPost, found, err := posts.GetById(ctx, h.FsClient, mediaRepy.VideoId)

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
	posts.Upsert(ctx, h.FsClient, exPost)

	return isHandeled, nil
}
