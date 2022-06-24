package messenger

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

func (m *Messenger) ReplyWithImage(chatId int64, replyToMessageId int, img []byte, imgName string, caption string) error {
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
