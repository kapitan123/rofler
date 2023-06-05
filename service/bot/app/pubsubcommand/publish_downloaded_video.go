package pubsubcommand

import (
	"context"
	"io"

	"github.com/kapitan123/telegrofler/service/bot/domain"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type PublishDownloadedVideo struct {
	messenger    messenger
	postsStorage postStorage
	filesBucket  filesBucket
}

type messenger interface {
	SendVideo(chatId domain.ChatId, videoId string, caption string, payload io.Reader) (int, error)
	ReplyWithText(chatId domain.ChatId, replyToMessageId domain.MessageId, text string) (int, error)
	Delete(chatId domain.ChatId, messageId domain.MessageId) error
}

type postStorage interface {
	GetByExternalSourceUrl(ctx context.Context, url string) (domain.Post, bool, error)
}

type filesBucket interface {
	Read(ctx context.Context, addr string, writer io.Writer) error
}

func NewPublishDownloadedVideo(messenger messenger, postsStorage postStorage, filesBucket filesBucket) *PublishDownloadedVideo {
	return &PublishDownloadedVideo{
		messenger:    messenger,
		postsStorage: postsStorage,
		filesBucket:  filesBucket,
	}
}

func (h *PublishDownloadedVideo) Handle(ctx context.Context, originalUrl string, savedAddr string) error {
	return nil
	post, found, err := h.postsStorage.GetByExternalSourceUrl(ctx, originalUrl)

	if err != nil {
		return err
	}

	if !found {
		return errors.Errorf("video %s was not found for url %s", savedAddr, originalUrl)
	}

	pr, pw := io.Pipe()

	errs := make(chan error, 1)

	go func() {
		defer pw.Close()
		errs <- h.filesBucket.Read(ctx, savedAddr, pw)
	}()

	_, err = h.messenger.SendVideo(post.ChatId, post.Id, post.Poster.AsUserMention(), pr)

	if err != nil {
		logrus.Error("can't post video to telegram channel", err)
	}

	if err := <-errs; err != nil {
		close(errs)
		return err
	}

	// it's really painful to perform integration tests as telegram has no backoff
	err = h.messenger.Delete(post.ChatId, post.OriginalMessageId)

	if err != nil {
		logrus.Error("can't delete message after publishing video", err)
	}

	return nil
}
