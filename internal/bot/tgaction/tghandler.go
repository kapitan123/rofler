package tgaction

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/data/firestore/posts"

	"github.com/kapitan123/telegrofler/internal/bot"
	log "github.com/sirupsen/logrus"
)

// Performs an action on recieving any message. Returns true if the execution should be stopped.
type BotMessageHandler interface {
	Handle(*tgbotapi.Message, context.Context) (bool, error)
}

// Executes the command from the message. Returns true if the command was handled.
type BotCommandHandler interface {
	BotMessageHandler
	GetCommandText() string
}

// Creates a map of handlers. Key is the command text, value is a handler.
func InitCommands(b *bot.Bot, ps *posts.PostsStorage) map[string]BotCommandHandler {
	commands := make(map[string]BotCommandHandler)
	commands[command1.GetCommandText()] = command1

	log.Infof("Commands registered %+T\n", commands)
	return commands
}

// Creates an array of handlers. Which will be tried to execute.
// Requires dependencies to be passed to the function
// Order of execution is determined by the order of the array.
func InitHandlers(b *bot.Bot, ps *posts.PostsStorage) *[]BotMessageHandler {
	// reply to 300 doesn't stp the execution of other handlers
	// Can handlers be jut functons with closures?
	handlers := []BotMessageHandler{
		NewReplyTo300(b),
		NewReplyToYes(b),
		NewReplyToNo(b),
		NewReplaceLinkWithMessage(b, ps),
		NewRecordBotPostReaction(b, ps),
		NewRecordReactionToUserMediaPost(b, ps),
	}

	log.Infof("Handlers registered %+T\n", handlers)
	return &handlers
}
