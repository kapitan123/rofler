package routes

import (
	"github.com/kapitan123/telegrofler/data/firestore"
	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/bot/tgaction"
)

type App struct {
	*bot.Bot
	*firestore.PostsStore
	handlers *[]tgaction.BotMessageHandler
	commands map[string]tgaction.BotCommandHandler
}

func NewApp() *App {
	//ctx := context.Background()
	// fsClient, err := firestore.NewClient(ctx, projectID)
	// if err != nil {
	// 	return nil, fmt.Errorf("firestore.NewClient: %v", err)
	// }

	return &App{
		Bot:        bot.New(),
		PostsStore: firestore.NewPostsStore(),
	}
}
