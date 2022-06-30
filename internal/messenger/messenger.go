package messenger

import (
	_ "embed"
	"errors"

	"github.com/kapitan123/telegrofler/internal/bot"

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

// AK TODO merge innternal calls to reduce dupliction
func (m *Messenger) SendImg(chatId int64, img []byte, imgName string, caption string) error {
	msg := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Name: imgName, Bytes: img})

	if caption != "" {
		msg.Caption = caption
	}

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

func (m *Messenger) GetUserCurrentProfilePic(userId int64) ([]byte, error) {
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

func (b *Messenger) Delete(chatId int64, messageId int) error {
	dmc := tgbotapi.DeleteMessageConfig{
		ChatID:    chatId,
		MessageID: messageId,
	}

	_, err := b.api.Request(dmc)
	if err != nil {
		return err
	}

	return nil
}

// Post TikTok back to the Telegram channel.
// Tags original poster and tiktok video info in description.
// remove unneeded depemdamcoes
// AK TODO rewrite without custom structure???
// separate this model from the extracted model
func (b *Messenger) SendTrackableVideo(tp *bot.SourceVideoPost) error {
	// Filename is id of the video
	fb := tgbotapi.FileBytes{Name: tp.VideoData.Id, Bytes: tp.VideoData.Payload}

	v := tgbotapi.NewVideo(tp.ChatId, fb)

	// AK TODO does it work with no duration?
	//v.Duration = tp.VideoData.Duration
	v.Caption = tp.GetCaption()
	v.ParseMode = tgbotapi.ModeHTML

	_, err := b.api.Send(v)

	return err
}

func (b *Messenger) GetAdminUserNames(chatId int64) ([]string, error) {
	req := tgbotapi.ChatAdministratorsConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID:             chatId,
			SuperGroupUsername: ""}, // AK TODO wtf is this parameter?
	}

	admins, err := b.api.GetChatAdministrators(req)

	if err != nil {
		return nil, err
	}

	userNames := make([]string, len(admins))

	for _, admin := range admins {
		userNames = append(userNames, admin.User.UserName)
	}

	return userNames, nil
}
