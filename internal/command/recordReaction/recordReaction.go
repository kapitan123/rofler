package recordReaction

import (
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/storage"
)

type RecordReaction struct {
	messenger messenger
	storage   postStorage
}

type (
	messenger interface {
		ReplyWithText(chatId int64, messageId int, text string) error
	}

	postStorage interface {
		GetPostById(ctx context.Context, videoId string) (storage.Post, bool, error)
		UpsertPost(ctx context.Context, p storage.Post) error
	}
)

func New(messenger messenger, storage postStorage) *RecordReaction {
	return &RecordReaction{
		messenger: messenger,
		storage:   storage,
	}
}

func (h *RecordReaction) Handle(ctx context.Context, m *tgbotapi.Message) error {
	mediaReply := extractUserMediaReaction(m)

	exPost, found, err := h.storage.GetPostById(ctx, mediaReply.VideoId)

	if err != nil {
		return err
	}

	from := m.ReplyToMessage.From

	if !found {
		reactions := make([]storage.Reaction, 0)
		roflerRef := storage.UserRef{
			Id:          from.ID,
			DisplayName: fmt.Sprintf("%s %s", from.FirstName, from.LastName),
		}

		// misc is video can be migrated to video type in the future
		postType := "pic"

		if m.ReplyToMessage.Video != nil {
			postType = "misc"
		}

		exPost = storage.Post{
			VideoId:   mediaReply.VideoId,
			Source:    postType,
			UserRef:   roflerRef,
			Url:       "",
			Reactions: reactions,
			PostedOn:  time.Now(),
			ChatId:    m.Chat.ID,
		}
	}

	exPost.AddReaction(mediaReply.ReactorRef, mediaReply.Text, mediaReply.ToMessageId)
	h.storage.UpsertPost(ctx, exPost)

	return nil
}

func (h *RecordReaction) ShouldRun(m *tgbotapi.Message) bool {
	rtm := m.ReplyToMessage

	if rtm == nil {
		return false
	}

	if m.From.ID == rtm.From.ID {
		return false
	}

	hasMedia := rtm.Video != nil || len(rtm.Photo) > 0

	if rtm == nil || !hasMedia || m.From.UserName == "" {
		return false
	}

	return true
}

type replyToMediaPost struct {
	VideoId     string
	ToMessageId int // RepllyToMessage.ID not the update.Message.ID
	ReactorRef  storage.UserRef
	Text        string
}

func extractUserMediaReaction(upd *tgbotapi.Message) replyToMediaPost {
	rtm := upd.ReplyToMessage

	// VideoId field also stores image id's should be converted with migration into another field
	mediaId := ""
	if rtm.Video != nil {
		mediaId = rtm.Video.FileID
	} else {
		mediaId = rtm.Photo[0].FileUniqueID
	}

	vr := replyToMediaPost{
		VideoId: mediaId,
		ReactorRef: storage.UserRef{
			Id:          rtm.From.ID,
			DisplayName: fmt.Sprintf("%s %s", upd.From.FirstName, upd.From.LastName),
		},
		ToMessageId: rtm.MessageID,
		Text:        upd.Text,
	}

	return vr
}
