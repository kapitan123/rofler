package bot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/kapitan123/telegrofler/config"
	"github.com/kapitan123/telegrofler/internal/data/post"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5" // https://go-telegram-bot-api.dev/
	log "github.com/sirupsen/logrus"
)

const (
	posterMaker = `🔥@(.*?)🔥`
)

// AK TODO compose an object from little mixin parts like TelegramPoster,
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

// Handles incoming chat messages.
// Tries to extract a TikTok video url from the message if no url was found returns nil
// Handles only mobile format
func (b *Bot) ConvertToSourceVideoPost(m *tgbotapi.Message) *SourceVideoPost {
	return &SourceVideoPost{
		Sender:            m.From.UserName,
		ChatId:            m.Chat.ID,
		Url:               m.Text,
		OriginalMessageId: m.MessageID,
	}
}

// AK TODO add sucess parameter
func (b *Bot) TryExtractVideoRepostReaction(upd *tgbotapi.Message) (post.VideoReaction, error) {
	vr := post.VideoReaction{}
	rtm := upd.ReplyToMessage

	if rtm == nil || rtm.From.UserName != "TelegroflBot" || rtm.Video == nil {
		return vr, nil
	}

	r := regexp.MustCompile(posterMaker)
	poster := r.FindStringSubmatch(rtm.Caption)[1]

	sender := upd.From.UserName
	// if the user reference his own post it is not a reaction. Can be moved outside the scope
	if sender == poster {
		return vr, nil
	}

	return post.VideoReaction{Sender: sender, VideoId: rtm.Video.FileName, Text: upd.Text, MessageId: upd.MessageID}, nil
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

func (b *Bot) GetCommand(m *tgbotapi.Message) (string, error) {
	// AK TODO handles only one command
	if !m.IsCommand() {
		return "", nil
	}

	if strings.HasPrefix("top", m.Command()) {
		return "top", nil
	}

	return "", nil
}

func (b *Bot) PostTopRofler(chatId int64, userName string, roflCount int) error {
	// AK TODO fetch rofler
	// AK TODO get user id
	// AK TODO fetch profile pic by Id
	// AK TODO apply watermark to image
	data := []byte("leleks")

	fb := tgbotapi.FileBytes{Name: "topRofler.jpg", Bytes: data}
	v := tgbotapi.NewPhoto(chatId, fb)

	v.Caption = fmt.Sprintf("\U0001F451 <b>@%s</b>\n<b>Likes:</b> %d", userName, roflCount)

	v.ParseMode = tgbotapi.ModeHTML

	_, err := b.api.Send(v)

	if err != nil {
		return err
	}

	return nil
}
