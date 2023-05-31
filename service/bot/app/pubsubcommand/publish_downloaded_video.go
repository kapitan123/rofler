package pubsubcommand

import (
	"context"
	"io"

	"github.com/kapitan123/telegrofler/service/bot/domain"
	"github.com/pkg/errors"
)

type PublishDownloadedVideo struct {
	messenger    messenger
	postsStorage postStorage
	filesBucket  filesBucket
}

type messenger interface {
	SendVideo(videoId string, chatId int64, caption string, payload io.Reader) (int, error)
}

type postStorage interface {
	GetByExternalSourceUrl(ctx context.Context, url string) (domain.Post, bool, error)
}

type filesBucket interface {
	Read(ctx context.Context, addr string) (io.Reader, error)
}

func NewPublishDownloadedVideo(messenger messenger, postsStorage postStorage, filesBucket filesBucket) *PublishDownloadedVideo {
	return &PublishDownloadedVideo{
		messenger:    messenger,
		postsStorage: postsStorage,
		filesBucket:  filesBucket,
	}
}

func (h *PublishDownloadedVideo) Handle(ctx context.Context, originalUrl string, savedAddr string) error {
	post, found, err := h.postsStorage.GetByExternalSourceUrl(ctx, originalUrl)

	if err != nil {
		return err
	}

	if !found {
		return errors.Errorf("video %s was not found for url %s", savedAddr, originalUrl)
	}

	reader, err := h.filesBucket.Read(ctx, savedAddr)

	if err != nil {
		return err
	}

	_, err = h.messenger.SendVideo(post.Id, post.ChatId, post.Poster.AsUserMention(), reader)

	if err != nil {
		return err
	}

	return err
}
