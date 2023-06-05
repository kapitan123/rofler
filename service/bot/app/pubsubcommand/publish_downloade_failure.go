package pubsubcommand

import (
	"context"

	"github.com/pkg/errors"
)

type PublishDownloadFailure struct {
	messenger    messenger
	postsStorage postStorage
}

func NewPublishDownloadFailure(messenger messenger, postsStorage postStorage) *PublishDownloadFailure {
	return &PublishDownloadFailure{
		messenger:    messenger,
		postsStorage: postsStorage,
	}
}

func (h *PublishDownloadFailure) Handle(ctx context.Context, originalUrl string) error {
	post, found, err := h.postsStorage.GetByExternalSourceUrl(ctx, originalUrl)

	if err != nil {
		return err
	}

	if !found {
		return errors.Errorf("video was not found for url %s", originalUrl)
	}

	_, err = h.messenger.ReplyWithText(post.ChatId, post.OriginalMessageId, "video failed to convert")

	return err
}
