package convertLinkToVideo

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/kapitan123/telegrofler/service/bot/domain"
	url "github.com/kapitan123/telegrofler/service/bot/domain/media"
	"github.com/kapitan123/telegrofler/service/bot/internal/messenger/format"
	"github.com/pkg/errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type RequestUrlContent struct {
	messenger    messenger
	postsStorage postStorage
	urlTopic     urlTopic
	urlsStorage  urlsStorage
}

type messenger interface {
	ReplyWithText(chatId int64, messageId int, text string) (int, error)
	SendVideo(chatId int64, trackToken string, caption string, payload io.Reader) (int, error)
	Delete(chatId int64, messageId int) error
}

type postStorage interface {
	UpsertPost(ctx context.Context, p domain.Post) error
}

type urlsStorage interface {
	CreateUrl(ctx context.Context, url string) error
}

type urlTopic interface {
	PublishUrl(ctx context.Context, url string) error
}

func New(messenger messenger, postsStorage postStorage, urlsStorage urlsStorage, urlTopic urlTopic) *RequestUrlContent {
	return &RequestUrlContent{
		messenger:    messenger,
		postsStorage: postsStorage,
		urlsStorage:  urlsStorage,
		urlTopic:     urlTopic,
	}
}

func (h *RequestUrlContent) Handle(ctx context.Context, m *tgbotapi.Message) error {
	url, chatId, senderId := m.Text, m.Chat.ID, m.From.ID
	senderName := fmt.Sprintf("%s %s", m.From.FirstName, m.From.LastName)

	err := h.urlsStorage.CreateUrl(ctx, url)

	if err != nil {
		return errors.Wrap(err, "unable to save found url")
	}

	err = h.urlTopic.PublishUrl(ctx, url)

	if err != nil {
		return errors.Wrap(err, "unable to publish found url to store")
	}

	err = h.messenger.Delete(chatId, m.MessageID)

	if err != nil {
		return errors.Wrap(err, "unable to delete message from chat")
	}

	log.Info("Url was found in a callback message: ", url)

	contentBuf := bytes.NewBuffer([]byte{})
	err = h.downloader.DownloadContent(meta.DownloadUrl, contentBuf)

	if err != nil {
		return err
	}

	// AK TODO get video from bucket
	// publish
	// delete firestore entry
	mention := format.AsUserMention(senderId, senderName)
	caption := fmt.Sprintf("<b>Rofler:</b> 🔥%s🔥\n<b>Title</b>: %s", mention, meta.Title)

	_, err = h.messenger.SendVideo(chatId, meta.Id, caption, contentBuf)

	if err != nil {
		return err
	}

	// we don't really care if if has failed and it makes integration tests a lot easier
	_ = h.messenger.Delete(chatId, m.MessageID)

	newPost := storage.Post{
		VideoId:   meta.Id,
		Source:    meta.Type,
		Url:       url,
		Reactions: []storage.Reaction{},
		PostedOn:  time.Now(),
	}

	err = h.postsStorage.UpsertPost(ctx, newPost)

	return err
}

func (h *RequestUrlContent) ShouldRun(m *tgbotapi.Message) bool {
	return len(m.Entities) > 0 && url.IsConvertable(m.Text)
}
