package tgaction

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/data/firestore/posts"
	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/source/sourceFactory"

	log "github.com/sirupsen/logrus"
)

type ReplaceLinkWithMessage struct {
	*bot.Bot
	postsStorage *posts.PostsStorage
}

func NewReplaceLinkWithMessage(b *bot.Bot, ps *posts.PostsStorage) *ReplaceLinkWithMessage {
	return &ReplaceLinkWithMessage{
		Bot:          b,
		postsStorage: ps,
	}
}

func (h *ReplaceLinkWithMessage) Handle(mess *tgbotapi.Message, ctx context.Context) (bool, error) {
	isHandeled := true
	extract, found := sourceFactory.TryGetExtractor(mess.Text)
	if !found {
		return !isHandeled, nil
	}

	svp := h.ConvertToSourceVideoPost(mess)

	if svp == nil {
		return !isHandeled, nil
	}

	log.Info("Url was found in a callback message: ", svp.Url)

	evi, err := extract(svp.Url)

	if err != nil {
		return isHandeled, err
	}

	svp.VideoData.Payload = evi.Payload
	svp.VideoData.Title = evi.Title
	svp.VideoData.Id = evi.Id

	log.Info("Trying to post to telegram: ", svp)

	err = h.RepostConvertedVideo(svp)

	if err != nil {
		return isHandeled, err
	}

	// we don't really care if if has failed and it makes integration tests a lot easier
	_ = h.DeletePost(svp.ChatId, svp.OriginalMessageId)

	newPost := posts.Post{
		VideoId:        svp.VideoData.Id,
		Source:         evi.Type,
		RoflerUserName: svp.Sender,
		Url:            svp.Url,
		Reactions:      []posts.Reaction{},
		KeyWords:       []string{},
		PostedOn:       time.Now(),
	}

	h.postsStorage.Upsert(ctx, newPost)

	if err != nil {
		return isHandeled, err
	}

	return isHandeled, nil
}
