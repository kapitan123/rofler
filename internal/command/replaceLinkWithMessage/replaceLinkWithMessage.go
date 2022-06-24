package replaceLinkWithMessage

import (
	"context"
	"time"

	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/source/sourceFactory"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/storage"
	log "github.com/sirupsen/logrus"
)

type ReplaceLinkWithMessage struct {
	messenger messenger
	storage   postStorage
}

type messenger interface {
	ReplyWithText(chatId int64, messageId int, text string) error
	RepostConvertedVideo(tp *bot.SourceVideoPost) error
}

type postStorage interface {
	GetById(ctx context.Context, videoId string) (storage.Post, bool, error)
	Upsert(ctx context.Context, p storage.Post) error
	Delete(chatId int64, messageId int) error
}

func New(messenger messenger) *ReplaceLinkWithMessage {
	return &ReplaceLinkWithMessage{
		messenger: messenger,
	}
}

func (h *ReplaceLinkWithMessage) Handle(ctx context.Context, m *tgbotapi.Message) error {
	// AK TODO extract for moking
	extract, _ := sourceFactory.TryGetExtractor(m.Text)

	svp := convertToSourceVideoPost(m)

	log.Info("Url was found in a callback message: ", svp.Url)

	evi, err := extract(svp.Url)

	if err != nil {
		return err
	}

	svp.VideoData.Payload = evi.Payload
	svp.VideoData.Title = evi.Title
	svp.VideoData.Id = evi.Id

	log.Info("Trying to post to telegram: ", svp)

	// AK TODO refactor to not use a custom structure?
	err = h.messenger.RepostConvertedVideo(svp)

	if err != nil {
		return err
	}

	// we don't really care if if has failed and it makes integration tests a lot easier
	_ = h.storage.Delete(svp.ChatId, svp.OriginalMessageId)

	newPost := storage.Post{
		VideoId:        svp.VideoData.Id,
		Source:         evi.Type,
		RoflerUserName: svp.Sender,
		Url:            svp.Url,
		Reactions:      []storage.Reaction{},
		KeyWords:       []string{},
		PostedOn:       time.Now(),
	}

	err = h.storage.Upsert(ctx, newPost)

	return err
}

func (h *ReplaceLinkWithMessage) ShouldRun(m *tgbotapi.Message) bool {
	// AK TODO extract for moking
	_, found := sourceFactory.TryGetExtractor(m.Text)
	return found
}

// Handles incoming chat messages.
// Tries to extract a TikTok video url from the message if no url was found returns nil
// Handles only mobile format
func convertToSourceVideoPost(m *tgbotapi.Message) *bot.SourceVideoPost {
	return &bot.SourceVideoPost{
		Sender:            m.From.UserName,
		ChatId:            m.Chat.ID,
		Url:               m.Text,
		OriginalMessageId: m.MessageID,
	}
}
