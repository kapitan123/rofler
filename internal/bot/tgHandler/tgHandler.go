package tghandler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// AK TODO come up with a way of execution so we don't do this tupid error handling
type BotCallbackHandler interface {
	Handle(*tgbotapi.Message) (bool, error)
}

// marks only commands
type BotCommand interface {
	*BotCallbackHandler
	GetCommandText() string
}
