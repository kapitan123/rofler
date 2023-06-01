package infra

import (
	_ "embed"
	"io"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/service/bot/domain"
	"github.com/sirupsen/logrus"
)

type TelegramMessenger struct {
	api *tgbotapi.BotAPI
}

func NewMessenger(token string) *TelegramMessenger {
	botapi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logrus.Panic("failed to create telegram client")
	}

	return &TelegramMessenger{
		api: botapi,
	}
}

func (m TelegramMessenger) SendText(chatID domain.ChatId, text string) (int, error) {
	msg := tgbotapi.NewMessage(chatID, text)

	msg.ParseMode = tgbotapi.ModeHTML

	res, err := m.api.Send(msg)

	if err != nil {
		return 0, err
	}

	return res.MessageID, err
}

func (m TelegramMessenger) ReplyWithText(chatId domain.ChatId, replyToMessageId int, caption string) (int, error) {
	msg := tgbotapi.NewMessage(chatId, caption)
	msg.ReplyToMessageID = replyToMessageId
	msg.ParseMode = tgbotapi.ModeHTML

	res, err := m.api.Send(msg)

	if err != nil {
		return 0, err
	}

	return res.MessageID, err
}

func (b TelegramMessenger) Delete(chatId domain.ChatId, messageId int) error {
	dmc := tgbotapi.DeleteMessageConfig{
		ChatID:    int64(chatId),
		MessageID: messageId,
	}

	_, err := b.api.Request(dmc)
	if err != nil {
		return err
	}

	return nil
}

func (b TelegramMessenger) SendVideo(videoId string, chatId domain.ChatId, caption string, payload io.Reader) (int, error) {
	vidbytes, err := io.ReadAll(payload)

	if err != nil {
		return 0, err
	}

	fb := tgbotapi.FileBytes{Name: videoId, Bytes: vidbytes}

	v := tgbotapi.NewVideo(int64(chatId), fb)

	v.Caption = caption
	v.ParseMode = tgbotapi.ModeHTML

	res, err := b.api.Send(v)

	if err != nil {
		return 0, err
	}

	return res.MessageID, err
}
