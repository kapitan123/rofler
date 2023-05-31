package tgcommand

import (
	"context"
	"fmt"
	"net/url"

	"github.com/kapitan123/telegrofler/common/logs"
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

func (h *RecordUrl) Handle(ctx context.Context, m message.Message) (err error) {
	defer func() {
		logs.LogExecutionResult(fmt.Sprintf("%T", h), m, err)
	}()

	url, err := m.GetEmbeddedUrl()

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

	// it's really painful to perform integration tests as telegram has no backoff
	_ = h.messenger.Delete(m.ChatId(), m.Id)

	return nil
}

func (h *RecordUrl) ShouldRun(m message.Message) bool {
	return m.HasDownloadableUrl()
}
