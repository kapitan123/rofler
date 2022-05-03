package routes

import (
	"github.com/kapitan123/telegrofler/data/firestore"
	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/bot/tgaction"
)

type API struct {
	*bot.Bot
	*firestore.PostsStore
	handlers *[]tgaction.BotMessageHandler
	commands map[string]tgaction.BotCommandHandler
}
