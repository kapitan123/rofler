package bot

import (
	"net/http"
	"regexp"

	"github.com/kapitan123/telegrofler/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5" // https://go-telegram-bot-api.dev/
	log "github.com/sirupsen/logrus"
)

const (
	mobilePrefixRegex = `https:\/\/[a-zA-Z]{2}\.tiktok\.com\/`
)

type Bot struct {
	api *tgbotapi.BotAPI
}

func New() (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(config.TelegramToken)
	if err != nil {
		log.Error("Error creating Telegram bot API: ", err)
		return nil, err
	}

	log.Info("Authorized on account %s", bot.Self.UserName)
	bot.Debug = true

	return &Bot{api: bot}, nil
}

// Post TikTok back to the Telegram channel.
// Tags original poster and tiktok video info in description.
func (b *Bot) PostTiktokVideoFromUrl(tp *TikTokVideoPost) error {
	// Filename is id of the video
	fb := tgbotapi.FileBytes{Name: tp.VideoData.Id, Bytes: tp.VideoData.Payload}

	v := tgbotapi.NewVideo(tp.ChatId, fb)

	v.Duration = tp.VideoData.Duration
	v.Caption = tp.GetCaption()
	v.ParseMode = tgbotapi.ModeHTML

	r, err := b.api.Send(v)

	log.Info("Message sent: ", r)

	if err != nil {
		return err
	}

	return nil
}

// Handles incoming chat messages.
// Tries to extract a TikTok video url from the message if no url was found returns nil
// Handles only mobile format
func (b *Bot) TryExtractTikTokUrlData(req *http.Request) (*TikTokVideoPost, error) {
	update, err := b.api.HandleUpdate(req)

	if err != nil {
		return nil, err
	}

	if r := regexp.MustCompile(mobilePrefixRegex); r.MatchString(update.Message.Text) {
		return &TikTokVideoPost{
			Sender:            update.Message.From.UserName,
			ChatId:            update.Message.Chat.ID,
			Url:               update.Message.Text,
			OriginalMessageId: update.Message.MessageID,
		}, nil
	}

	return nil, nil
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
