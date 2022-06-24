package messenger

import (
	_ "embed"
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/source"
)

type Messenger struct {
	api *tgbotapi.BotAPI
}

func New(api *tgbotapi.BotAPI) *Messenger {
	return &Messenger{api: api}
}

func (m *Messenger) SendText(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := m.api.Send(msg)
	return err
}

func (m *Messenger) ReplyWithText(chatId int64, replyToMessageId int, caption string) error {
	msg := tgbotapi.NewMessage(chatId, caption)
	msg.ReplyToMessageID = replyToMessageId

	_, err := m.api.Send(msg)

	if err != nil {
		return err
	}

	return nil
}

func (m *Messenger) ReplyWithImg(chatId int64, replyToMessageId int, img []byte, imgName string, caption string) error {
	msg := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Name: imgName, Bytes: img})
	msg.ReplyToMessageID = replyToMessageId

	if caption != "" {
		msg.Caption = caption
	}

	_, err := m.api.Send(msg)

	if err != nil {
		return err
	}

	return nil
}

func (m *Messenger) GetCurrentUserProfilePic(userId int64) ([]byte, error) {
	ppicReq := tgbotapi.UserProfilePhotosConfig{
		UserID: userId,
		Offset: 0,
		Limit:  1,
	}

	pics, err := m.api.GetUserProfilePhotos(ppicReq)
	if err != nil {
		return nil, err
	}

	if len(pics.Photos) == 0 {
		return nil, errors.New("no profile picture was found")
	}

	ppicMeta := pics.Photos[0][2]

	ppic, err := m.api.GetFile(tgbotapi.FileConfig{FileID: ppicMeta.FileID})

	if err != nil {
		return nil, err
	}

	downloadLink := ppic.Link(m.api.Token)

	// AK TODO this crap is super slow
	return source.DownloadBytesFromUrl(downloadLink)
}

func (m *Messenger) Delete(chatId int64, messageId int) error {
	dmc := tgbotapi.DeleteMessageConfig{
		ChatID:    chatId,
		MessageID: messageId,
	}

	_, err := m.api.Request(dmc)
	if err != nil {
		return err
	}

	return nil
}
