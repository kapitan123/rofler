package replaceLinkWithMessage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/kapitan123/telegrofler/internal/contentLoader"
	"github.com/kapitan123/telegrofler/internal/messenger/format"
	"github.com/kapitan123/telegrofler/internal/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type ReplaceLinkWithMessage struct {
	messenger  messenger
	storage    postStorage
	downloader downloader
}

type messenger interface {
	ReplyWithText(chatId int64, messageId int, text string) error
	SendVideo(chatId int64, trackToken string, caption string, payload io.Reader) error
	Delete(chatId int64, messageId int) error
}

type postStorage interface {
	UpsertPost(ctx context.Context, p storage.Post) error
}

type downloader interface {
	DownloadContent(dUrl string, res io.Writer) error
	ExtractVideoMeta(url string) (*contentLoader.VideoMeta, error)
	CanExtractVideoMeta(url string) bool
}

func New(messenger messenger, storage postStorage, downloader downloader) *ReplaceLinkWithMessage {
	return &ReplaceLinkWithMessage{
		messenger:  messenger,
		storage:    storage,
		downloader: downloader,
	}
}

func (h *ReplaceLinkWithMessage) Handle(ctx context.Context, m *tgbotapi.Message) error {
	url, chatId, senderId := m.Text, m.Chat.ID, m.From.ID
	senderName := fmt.Sprintf("%s %s", m.From.FirstName, m.From.LastName)

	meta, err := h.downloader.ExtractVideoMeta(url)

	if err != nil {
		return err
	}

	log.Info("Url was found in a callback message: ", url)

	contentBuf := bytes.NewBuffer([]byte{})
	err = h.downloader.DownloadContent(meta.DownloadUrl, contentBuf)

	if err != nil {
		return err
	}

	mention := format.AsUserMention(senderId, senderName)
	caption := fmt.Sprintf("<b>Rofler:</b> 🔥%s🔥\n<b>Title</b>: %s", mention, meta.Title)

	err = h.messenger.SendVideo(chatId, meta.Id, caption, contentBuf)

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

	err = h.storage.UpsertPost(ctx, newPost)

	return err
}

func (h *ReplaceLinkWithMessage) ShouldRun(m *tgbotapi.Message) bool {
	return h.downloader.CanExtractVideoMeta(m.Text)
}
