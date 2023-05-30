package recordBotPostReaction

import (
	"context"
	"fmt"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/service/bot/internal/storage"
)

var posterMakerRegExp = regexp.MustCompile(`🔥(.*?)🔥`)

type RecordBotPostReaction struct {
	messenger messenger
	storage   postStorage
}

type (
	messenger interface {
		ReplyWithText(chatId int64, messageId int, text string) (int, error)
	}

	postStorage interface {
		GetPostById(ctx context.Context, videoId string) (storage.Post, bool, error)
		UpsertPost(ctx context.Context, p storage.Post) error
	}
)

func New(messenger messenger, storage postStorage) *RecordBotPostReaction {
	return &RecordBotPostReaction{
		messenger: messenger,
		storage:   storage,
	}
}

func (h *RecordBotPostReaction) Handle(ctx context.Context, m *tgbotapi.Message) error {
	reply := extractVideoRepostReaction(m)

	exPost, found, err := h.storage.GetPostById(ctx, reply.VideoId)

	if err != nil {
		return err
	}

	// in this case we don't record reaction as all bot posts should be saved already
	if !found {
		return nil
	}

	exPost.AddReaction(reply.ReactorRef, reply.Text, reply.ToMessageId)
	h.storage.UpsertPost(ctx, exPost)
	return nil
}

func (h *RecordBotPostReaction) ShouldRun(m *tgbotapi.Message) bool {
	return containsVideoRepostReaction(m)
}

// AK TODO this is part of the domain
type replyToMediaPost struct {
	VideoId     string
	ToMessageId int // RepllyToMessage.ID not the update.Message.ID
	ReactorRef  storage.UserRef
	Text        string
}

func extractVideoRepostReaction(upd *tgbotapi.Message) replyToMediaPost {
	rtm := upd.ReplyToMessage
	reply := replyToMediaPost{
		VideoId: rtm.Video.FileName,
		ReactorRef: storage.UserRef{
			Id:          rtm.From.ID,
			DisplayName: fmt.Sprintf("%s %s", upd.From.FirstName, upd.From.LastName),
		},
		Text:        upd.Text,
		ToMessageId: rtm.MessageID,
	}

	return reply
}

func containsVideoRepostReaction(upd *tgbotapi.Message) bool {
	rtm := upd.ReplyToMessage

	if rtm == nil || rtm.From.UserName != "TelegroflBot" || rtm.Video == nil {
		return false
	}

	// monitored bot posts has exactly one mention
	if !posterMakerRegExp.MatchString(rtm.Caption) {
		return false
	}

	for _, e := range rtm.CaptionEntities {
		// only if reactor is not a mentioned in caption user (not an orignial poster)
		if e.Type == "text_mention" && e.User.ID != upd.From.ID {
			return true
		}
	}

	return false
}
