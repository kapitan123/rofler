package routes

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/kapitan123/telegrofler/config"
	"github.com/kapitan123/telegrofler/data/firestore/posts"
	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/bot/tgaction"
)

type App struct {
	*bot.Bot
	PostsStorage *posts.PostsStorage
	handlers     *[]tgaction.BotMessageHandler
	commands     map[string]tgaction.BotCommandHandler
}

func NewApp() (*App, error) {
	ctx := context.Background()
	fsClient, err := firestore.NewClient(ctx, config.ProjectId)
	if err != nil {
		return nil, fmt.Errorf("firestore.NewClient: %v", err)
	}

	return &App{
		Bot:          bot.New(),
		PostsStorage: posts.New(fsClient),
	}, nil
}

func (app *App) Close() {
	app.PostsStorage.Close()
}
