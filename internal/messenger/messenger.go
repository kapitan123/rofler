package messenger

import (
	_ "embed"
	"errors"
	"io"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// AK TODO convert all byte arrays to io readers and writers
type Messenger struct {
	api        *tgbotapi.BotAPI
	downloader downloader
}

type downloader interface {
	DownloadContent(dUrl string, res io.Writer) error
}

func New(api *tgbotapi.BotAPI, downloader downloader) *Messenger {
	return &Messenger{
		api:        api,
		downloader: downloader,
	}
}

func (m *Messenger) SendText(chatID int64, text string) (int, error) {
	msg := tgbotapi.NewMessage(chatID, text)

	msg.ParseMode = tgbotapi.ModeHTML

	res, err := m.api.Send(msg)

	if err != nil {
		return 0, err
	}

	return res.MessageID, err
}

func (m *Messenger) ReplyWithText(chatId int64, replyToMessageId int, caption string) (int, error) {
	msg := tgbotapi.NewMessage(chatId, caption)
	msg.ReplyToMessageID = replyToMessageId
	msg.ParseMode = tgbotapi.ModeHTML

	res, err := m.api.Send(msg)

	if err != nil {
		return 0, err
	}

	return res.MessageID, err
}

// AK TODO merge innternal calls to reduce dupliction
func (m *Messenger) SendImg(chatId int64, img io.Reader, imgName string, caption string) (int, error) {
	imgBytes, err := io.ReadAll(img)

	if err != nil {
		return 0, err
	}

	msg := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Name: imgName, Bytes: imgBytes})
	msg.ParseMode = tgbotapi.ModeHTML
	if caption != "" {
		msg.Caption = caption
	}

	res, err := m.api.Send(msg)

	if err != nil {
		return 0, err
	}

	return res.MessageID, err
}

func (m *Messenger) ReplyWithImg(chatId int64, replyToMessageId int, img io.Reader, imgName string, caption string) (int, error) {
	imgbytes, err := io.ReadAll(img)
	if err != nil {
		return 0, err
	}

	msg := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Name: imgName, Bytes: imgbytes})
	msg.ReplyToMessageID = replyToMessageId
	msg.ParseMode = tgbotapi.ModeHTML

	if caption != "" {
		msg.Caption = caption
	}

	_, err = m.api.Send(msg)

	res, err := m.api.Send(msg)

	if err != nil {
		return 0, err
	}

	return res.MessageID, err
}

func (m *Messenger) GetUserCurrentProfilePic(userId int64, res io.Writer) error {
	ppicReq := tgbotapi.UserProfilePhotosConfig{
		UserID: userId,
		Offset: 0,
		Limit:  1,
	}

	pics, err := m.api.GetUserProfilePhotos(ppicReq)
	if err != nil {
		return err
	}

	if len(pics.Photos) == 0 {
		return errors.New("no profile picture was found")
	}

	ppicMeta := pics.Photos[0][2]

	ppic, err := m.api.GetFile(tgbotapi.FileConfig{FileID: ppicMeta.FileID})

	if err != nil {
		return err
	}

	downloadLink := ppic.Link(m.api.Token)

	// AK TODO this crap is super slow
	return m.downloader.DownloadContent(downloadLink, res)
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

func (b *Messenger) SendVideo(chatId int64, videoId string, caption string, payload io.Reader) (int, error) {
	vidbytes, err := io.ReadAll(payload)

	if err != nil {
		return 0, err
	}
	// Filename is id of the video
	fb := tgbotapi.FileBytes{Name: videoId, Bytes: vidbytes}

	v := tgbotapi.NewVideo(chatId, fb)

	v.Caption = caption
	v.ParseMode = tgbotapi.ModeHTML

	res, err := b.api.Send(v)

	if err != nil {
		return 0, err
	}

	return res.MessageID, err
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

type UserRef struct {
	Id        string
	UserName  string
	UserTitle string
}
