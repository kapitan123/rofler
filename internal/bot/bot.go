package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	_ "embed"

	"github.com/kapitan123/telegrofler/config"
	"github.com/kapitan123/telegrofler/pkg/source"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5" // https://go-telegram-bot-api.dev/
	log "github.com/sirupsen/logrus"
)

// A proxy for the Telegram API.
// Encapsulates all logic related to telegram interactions
type Bot struct {
	api *tgbotapi.BotAPI
}

func New() *Bot {
	bot, err := tgbotapi.NewBotAPI(config.TelegramToken)
	if err != nil {
		log.Panic("Error creating Telegram bot API: ", err)
	}

	log.Info("Authorized on account ", bot.Self.UserName)
	bot.Debug = true

	return &Bot{api: bot}
}

// Post TikTok back to the Telegram channel.
// Tags original poster and tiktok video info in description.
func (b *Bot) RepostConvertedVideo(tp *SourceVideoPost) error {
	// Filename is id of the video
	fb := tgbotapi.FileBytes{Name: tp.VideoData.Id, Bytes: tp.VideoData.Payload}

	v := tgbotapi.NewVideo(tp.ChatId, fb)

	// AK TODO does it work with no duration?
	//v.Duration = tp.VideoData.Duration
	v.Caption = tp.GetCaption()
	v.ParseMode = tgbotapi.ModeHTML

	r, err := b.api.Send(v)

	log.Info("Message sent: ", r)

	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) DeletePost(chatId int64, messageId int) error {
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

func (b *Bot) GetUpdate(req *http.Request) (*tgbotapi.Update, error) {
	update, err := b.api.HandleUpdate(req)
	if err != nil {
		return nil, err
	}
	ujs, err := json.Marshal(update)
	log.Info("Callback content:", string(ujs))

	if err != nil {
		return nil, err
	}

	return update, nil
}

// AK TODO remove after migration to a more abstract wrapper
func (b *Bot) PostTopRofler(chatId int64, userName string, roflCount int) error {
	// AK TODO fetch rofler
	// AK TODO get user id
	// AK TODO fetch profile pic by Id
	// AK TODO apply watermark to image
	//data := []byte("leleks")

	//fb := tgbotapi.FileBytes{Name: "topRofler.jpg", Bytes: data}
	//v := tgbotapi.NewPhoto(chatId, fb)

	//v.Caption = fmt.Sprintf("\U0001F451 <b>@%s</b>\n<b>Likes:</b> %d", userName, roflCount)
	v := tgbotapi.NewMessage(chatId, fmt.Sprintf("\U0001F451 <b>@%s</b>\n<b>Likes:</b> %d", userName, roflCount))
	v.ParseMode = tgbotapi.ModeHTML

	_, err := b.api.Send(v)

	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) PostReplyWithText(chatId int64, replyToMessageId int, caption string) error {
	msg := tgbotapi.NewMessage(chatId, caption)
	msg.ReplyToMessageID = replyToMessageId

	_, err := b.api.Send(msg)

	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) GetCurrentUserProfilePic(userId int64) ([]byte, error) {
	ppicReq := tgbotapi.UserProfilePhotosConfig{
		UserID: userId,
		Offset: 0,
		Limit:  1,
	}

	pics, err := b.api.GetUserProfilePhotos(ppicReq)
	if err != nil {
		return nil, err
	}

	if len(pics.Photos) == 0 {
		return nil, errors.New("no profile picture was found")
	}

	ppicMeta := pics.Photos[0][2]

	ppic, err := b.api.GetFile(tgbotapi.FileConfig{FileID: ppicMeta.FileID})

	if err != nil {
		return nil, err
	}

	downloadLink := ppic.Link(b.api.Token)

	return source.DownloadBytesFromUrl(downloadLink)
}

func (b *Bot) PostReplyWithImage(chatId int64, replyToMessageId int, img []byte, imgName string, caption string) error {
	msg := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Name: imgName, Bytes: img})
	msg.ReplyToMessageID = replyToMessageId

	if caption != "" {
		msg.Caption = caption
	}

	_, err := b.api.Send(msg)

	if err != nil {
		return err
	}

	return nil
}
