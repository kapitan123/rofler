package recordBotPostReaction

import (
	"context"
	"fmt"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/service/bot/domain"
	"github.com/kapitan123/telegrofler/service/bot/domain/message"
)

var posterMakerRegExp = regexp.MustCompile(`🔥(.*?)🔥`)

type RecordBotPostReaction struct {
	storage postStorage
}

type (
	postStorage interface {
		GetPostById(ctx context.Context, videoId string) (domain.Post, bool, error)
		UpsertPost(ctx context.Context, p domain.Post) error
	}
)

func New(storage postStorage) *RecordBotPostReaction {
	return &RecordBotPostReaction{
		storage: storage,
	}
}

func (h *RecordBotPostReaction) Handle(ctx context.Context, m message.Message) error {
	existingPost, found, err := h.storage.GetPostById(ctx, m.MediaId())

	if err != nil {
		return err
	}

	// in this case we don't record reaction as all bot posts should be saved already
	if !found {
		return fmt.Errorf("original post was not saved, reaction will not be recorded")
	}

	existingPost.AddReaction(m.AsReaction())

	h.storage.UpsertPost(ctx, existingPost)
	return nil
}

func (h *RecordBotPostReaction) ShouldRun(m *tgbotapi.Message) bool {
	return isReplyToVideoPostedByBot(m)
}

func isReplyToVideoPostedByBot(upd *tgbotapi.Message) bool {
	rtm := upd.ReplyToMessage

	if rtm == nil || rtm.Video == nil {
		return false
	}

	isVideoConvertedByBot := posterMakerRegExp.MatchString(rtm.Caption) || rtm.From.UserName == "TelegroflBot"

	if !isVideoConvertedByBot {
		return false
	}

	// AK TODO all this logic should belong to domain
	// AKA tg message. or should be just an extension of tgbotApi

	for _, e := range rtm.CaptionEntities {
		// only if reactor is not mentioned in caption user (not an orignial poster)
		if e.Type == "text_mention" && e.User.ID != upd.From.ID {
			return true
		}
	}

	return false
}
