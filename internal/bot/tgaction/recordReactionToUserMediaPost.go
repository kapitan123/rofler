package tgaction

import (
	"time"

	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/data/post"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type RecordReactionToUserMediaPost struct {
	*bot.Bot
	*post.PostsStore
}

func NewRecordReactionToUserMediaPost(b *bot.Bot, ps *post.PostsStore) *RecordReactionToUserMediaPost {
	return &RecordReactionToUserMediaPost{
		Bot:        b,
		PostsStore: ps,
	}
}

func (h *RecordReactionToUserMediaPost) Handle(m *tgbotapi.Message) (bool, error) {
	rtm := m.ReplyToMessage

	if rtm == nil || rtm.Video == nil {
		return false, nil
	}

	isHandeled := true

	reaction, err := bot.ExtractUserMediaReaction(m)

	if err != nil {
		return !isHandeled, err
	}

	// AK TODO should actually return nil
	if reaction.Sender == "" {
		return !isHandeled, nil
	}

	log.Infof("Reaction was found for %s sent by %s", reaction.VideoId, reaction.Sender)

	exPost, found, err := h.GetById(reaction.VideoId)

	if err != nil {
		return isHandeled, err
	}

	if !found {
		reactions := make([]post.Reaction, 0)
		exPost = post.Post{
			VideoId:        reaction.VideoId,
			Source:         "misc",
			RoflerUserName: rtm.From.UserName,
			Url:            "",
			Reactions:      reactions,
			PostedOn:       time.Now(),
		}
	}

	exPost.AddReaction(reaction.Sender, reaction.Text, reaction.MessageId)
	h.Upsert(exPost)

	return isHandeled, nil
}
