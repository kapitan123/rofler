package toprofl

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/storage"
)

const commandName = "toprofl"

type TopRofl struct {
	messenger messenger
	storage   postStorage
}

type messenger interface {
	ReplyWithText(chatId int64, replyToMessageId int, caption string) error
	SendText(chatId int64, text string) error
	GetChatAdmins(chatId int64) ([]tgbotapi.ChatMember, error)
}

type postStorage interface {
	GetLastWeekPosts(ctx context.Context, chatid int64) ([]storage.Post, error)
}

func New(messenger messenger, storage postStorage) *TopRofl {
	return &TopRofl{
		messenger: messenger,
		storage:   storage,
	}
}

func (h *TopRofl) Handle(ctx context.Context, message *tgbotapi.Message) error {
	posts, err := h.storage.GetLastWeekPosts(ctx, message.Chat.ID)
	if err != nil {
		return err
	}

	topRofl := findMax(posts)
	reactionCount := len(topRofl.Reactions)

	if reactionCount == 0 {
		err = h.messenger.SendText(message.Chat.ID, "На этой неделе смешного не было")
	} else {
		mesId := topRofl.Reactions[0].MessageId
		cap := fmt.Sprintf("Самый смешный пост: %d 👍", reactionCount)
		err = h.messenger.ReplyWithText(message.Chat.ID, mesId, cap)
	}
	if err != nil {
		return err
	}
	return nil
}

func findMax(posts []storage.Post) storage.Post {
	max := storage.Post{
		Reactions: []storage.Reaction{},
	}

	for _, p := range posts {
		if len(max.Reactions) < len(p.Reactions) {
			max = p
		}
	}
	return max
}

func (h *TopRofl) ShouldRun(message *tgbotapi.Message) bool {
	return message.IsCommand() && message.Command() == commandName
}
