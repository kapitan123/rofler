package tghandler

import (
	log "github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/data/post"
)

// AK TODO come up with a way of execution so we don't do this tupid error handling
type BotMessageHandler interface {
	Handle(*tgbotapi.Message) (bool, error)
}

// marks only commands
type BotCommandHandler interface {
	BotMessageHandler
	GetCommandText() string
}

func InitCommands(b *bot.Bot, ps *post.PostsStore) map[string]BotCommandHandler {
	commands := make(map[string]BotCommandHandler)
	command1 := NewPostTopRoflerCommand(b, ps)
	commands[command1.GetCommandText()] = command1
	return commands
}

func InitHandlers(b *bot.Bot, ps *post.PostsStore) []BotMessageHandler {
	hadler1 := NewPostTopRoflerCommand(b, ps)
	hadler2 := NewRecordReaction(b, ps)
	hadler3 := NewReplaceLinkWithMessage(b, ps)
	handlers := []BotMessageHandler{hadler1, hadler2, hadler3}

	log.Infof("%+v\n", handlers)
	return handlers
}
