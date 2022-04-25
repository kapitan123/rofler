package tgHandler

import (
	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/data"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PostTopRofler struct {
	*bot.Bot
	*data.PostsStore
}

func (h *PostTopRofler) Handle(m *tgbotapi.Message) (bool, error) {
	tr, roflCount, err := h.GetTopRoflerFromPosts()
	if err != nil {
		return false, err
	}

	err = h.PostTopRofler(m.Chat.ID, tr, roflCount)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (h *PostTopRofler) GetCommandText() string {
	return "top"
}
