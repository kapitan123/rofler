package tgcommand

import (
	"context"
	"net/url"

	"github.com/kapitan123/telegrofler/service/bot/domain"
	"github.com/kapitan123/telegrofler/service/bot/domain/message"
	"github.com/pkg/errors"
)

type RecordUrl struct {
	messenger    messenger
	postsStorage postStorage
	urlTopic     urlTopic
}

type urlTopic interface {
	PublishUrl(ctx context.Context, url *url.URL) error
}

func NewRecordUrl(messenger messenger, postsStorage postStorage, urlTopic urlTopic) *RecordUrl {
	return &RecordUrl{
		messenger:    messenger,
		postsStorage: postsStorage,
		urlTopic:     urlTopic,
	}
}

func (h *RecordUrl) Handle(ctx context.Context, m message.Message) error {
	url, _, err := m.FindUrl()

	if err != nil {
		return err
	}

	post := domain.NewPostFromExternalSource(url, m.From(), m.ChatId())

	err = h.postsStorage.UpsertPost(ctx, post)

	if err != nil {
		return errors.Wrap(err, "unable to save found url")
	}

	err = h.urlTopic.PublishUrl(ctx, url)

	if err != nil {
		return errors.Wrap(err, "unable to publish found url to pubsub")
	}

	err = h.messenger.Delete(m.ChatId(), m.Id)

	if err != nil {
		return errors.Wrap(err, "unable to delete message from chat")
	}

	return nil
}

func (h *RecordUrl) ShouldRun(m message.Message) bool {
	return m.HasDownloadableUrl()
}
