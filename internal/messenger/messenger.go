package messenger

import (
	_ "embed"
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Messenger struct {
	api        *tgbotapi.BotAPI
	downloader downloader
}

type downloader interface {
	DownloadContent(dUrl string) ([]byte, error)
}

func New(api *tgbotapi.BotAPI, downloader downloader) *Messenger {
	return &Messenger{
		api:        api,
		downloader: downloader,
	}
}

func (m *Messenger) SendText(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)

	msg.ParseMode = tgbotapi.ModeHTML

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
	return m.downloader.DownloadContent(downloadLink)
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
type VideoData struct {
	Id      string
	Title   string
	Payload []byte
}

func (b *Messenger) SendTrackableVideo(chatId int64, linktToUserName string, trackToken string, title string, payload []byte) error {
	// Filename is id of the video
	fb := tgbotapi.FileBytes{Name: trackToken, Bytes: payload}

	v := tgbotapi.NewVideo(chatId, fb)

	v.Caption = GetTrackableCaption(linktToUserName, title)
	v.ParseMode = tgbotapi.ModeHTML

	_, err := b.api.Send(v)

	return err
}

func GetTrackableCaption(linktToUserName string, title string) string {
	return fmt.Sprintf("<b>Rofler:</b> 🔥@%s🔥\n<b>Title</b>: %s", linktToUserName, title)
}

// AK we can abstract away ChatMember. But all the commands already has a dependency on tgBotApi.
func (b *Messenger) GetChatAdmins(chatId int64) ([]tgbotapi.ChatMember, error) {
	req := tgbotapi.ChatAdministratorsConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID:             chatId,
			SuperGroupUsername: ""}, // AK TODO wtf is this parameter?
	}

	// doc claims that there will be no bots in results but it is not true
	admins, err := b.api.GetChatAdministrators(req)

	if err != nil {
		return nil, err
	}

	noBots := make([]tgbotapi.ChatMember, 0)
	for _, a := range admins {
		if !a.User.IsBot {
			noBots = append(noBots, a)
		}
	}

	return noBots, nil
}
