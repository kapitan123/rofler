package replaceLinkWithMessage

import (
	"context"
	"time"

	"github.com/kapitan123/telegrofler/internal/command/replaceLinkWithMessage/source/sourceFactory"

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
	SendTrackableVideo(chatId int64, linktToUserName string, trackToken string, title string, payload []byte) error
	Delete(chatId int64, messageId int) error
}

type postStorage interface {
	GetById(ctx context.Context, videoId string) (storage.Post, bool, error)
	UpsertPost(ctx context.Context, p storage.Post) error
}

func New(messenger messenger, storage postStorage) *ReplaceLinkWithMessage {
	return &ReplaceLinkWithMessage{
		messenger: messenger,
		storage:   storage,
	}
}

func (h *ReplaceLinkWithMessage) Handle(ctx context.Context, m *tgbotapi.Message) error {
	// AK TODO extract for moking
	// video downloads are specific to this handler, so ex
	extract, _ := sourceFactory.TryGetExtractor(m.Text)

	url, chatId, sender := m.Text, m.Chat.ID, m.From.UserName

	log.Info("Url was found in a callback message: ", url)

	evi, err := extract(url)

	if err != nil {
		return err
	}

	err = h.messenger.SendTrackableVideo(chatId, sender, evi.Id, evi.Title, evi.Payload)

	if err != nil {
		return err
	}

	// we don't really care if if has failed and it makes integration tests a lot easier
	_ = h.messenger.Delete(chatId, m.MessageID)

	newPost := storage.Post{
		VideoId:        evi.Id,
		Source:         evi.Type,
		RoflerUserName: sender,
		Url:            url,
		Reactions:      []storage.Reaction{},
		KeyWords:       []string{},
		PostedOn:       time.Now(),
	}

	err = h.storage.UpsertPost(ctx, newPost)

	return err
}

func (h *ReplaceLinkWithMessage) ShouldRun(m *tgbotapi.Message) bool {
	// AK TODO extract for moking
	// separate extractor and cheking of it existance
	_, found := sourceFactory.TryGetExtractor(m.Text)
	return found
}
