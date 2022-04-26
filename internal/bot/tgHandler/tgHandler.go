package tghandler

import (
	log "github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/data/post"
)

// Performs an action on recieving any message. Returns true if the execution should be stopped.
type BotMessageHandler interface {
	Handle(*tgbotapi.Message) (bool, error)
}

// Executes the command from the message. Returns true if the command was handled.
type BotCommandHandler interface {
	BotMessageHandler
	GetCommandText() string
}

// Creates a map of handlers. Key is the command text, value is a handler.
func InitCommands(b *bot.Bot, ps *post.PostsStore) map[string]BotCommandHandler {
	commands := make(map[string]BotCommandHandler)
	command1 := NewPostTopRoflerCommand(b, ps)
	commands[command1.GetCommandText()] = command1
	return commands
}

// Creates an array of handlers. Which will be tried to execute.
// Requires dependencies to be passed to the function
// Order of execution is determined by the order of the array.
func InitHandlers(b *bot.Bot, ps *post.PostsStore) []BotMessageHandler {
	hadler1 := NewPostTopRoflerCommand(b, ps)
	hadler2 := NewRecordReaction(b, ps)
	hadler3 := NewReplaceLinkWithMessage(b, ps)
	handlers := []BotMessageHandler{hadler1, hadler2, hadler3}

	log.Infof("%+v\n", handlers)
	return handlers
}
