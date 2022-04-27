package routes

import (
	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/bot/tgaction"
	"github.com/kapitan123/telegrofler/internal/data/post"
)

// AK TODO should be in a separate file

type API struct {
	*bot.Bot
	*post.PostsStore
	handlers *[]tgaction.BotMessageHandler
	commands map[string]tgaction.BotCommandHandler
}
