package convertLinkToVideo

import (
	"context"
	"fmt"
	"io"

	"github.com/kapitan123/telegrofler/service/bot/domain"
	"github.com/kapitan123/telegrofler/service/bot/domain/media"
)

type RequestUrlContent struct {
	messenger    messenger
	postsStorage postStorage
	urlsStorage  mediaStorage
	filesBucket  filesBucket
}

type messenger interface {
	SendVideo(chatId int64, caption string, payload io.Reader) (int, error)
}

type postStorage interface {
	UpsertPost(ctx context.Context, p domain.Post) error
}

type mediaStorage interface {
	GetByUrl(ctx context.Context, url string) (media.Media, bool, error)
}

type filesBucket interface {
	Read(ctx context.Context, addr string) (io.Reader, error)
}

func New(messenger messenger, postsStorage postStorage, urlsStorage mediaStorage, filesBucket filesBucket) *RequestUrlContent {
	return &RequestUrlContent{
		messenger:    messenger,
		postsStorage: postsStorage,
		urlsStorage:  urlsStorage,
		filesBucket:  filesBucket,
	}
}

func (h *RequestUrlContent) Handle(ctx context.Context, originalUrl string, savedAddr string) error {
	// AK TODO get video from bucket
	// publish
	// delete firestore entry
	media, found, err := h.urlsStorage.GetByUrl(ctx, originalUrl)

	if err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("video %s was not found for url %s", savedAddr, originalUrl)
	}

	reader, err := h.filesBucket.Read(ctx, savedAddr)

	if err != nil {
		return err
	}

	// AK TODO might be moved to domain
	caption := fmt.Sprintf("<b>Rofler:</b> 🔥%s🔥", media.Poster.AsUserMention())

	_, err = h.messenger.SendVideo(media.ChatId, caption, reader)

	if err != nil {
		return err
	}

	// AK TODO Do I even need a new collection? looks like media can be removed or renamed to posts
	newPost := domain.NewPost(originalUrl, media.Poster, media.ChatId)

	err = h.postsStorage.UpsertPost(ctx, newPost)

	return err
}
