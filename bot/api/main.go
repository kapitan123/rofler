package main

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/kapitan123/telegrofler/bot/internal/command"
	choosePidor "github.com/kapitan123/telegrofler/bot/internal/command/choose_pidor"
	recordBotPostReaction "github.com/kapitan123/telegrofler/bot/internal/command/record_bot_post_reaction"
	recordReaction "github.com/kapitan123/telegrofler/bot/internal/command/record_reaction"
	replaceLinkWithMessage "github.com/kapitan123/telegrofler/bot/internal/command/replace_link_with_message"
	replyTo300 "github.com/kapitan123/telegrofler/bot/internal/command/reply_to_300"
	replyToNo "github.com/kapitan123/telegrofler/bot/internal/command/reply_to_no"
	replyToYes "github.com/kapitan123/telegrofler/bot/internal/command/reply_to_yes"
	toppidor "github.com/kapitan123/telegrofler/bot/internal/command/top_pidor"
	toprofler "github.com/kapitan123/telegrofler/bot/internal/command/top_rofler"
	contentLoader "github.com/kapitan123/telegrofler/bot/internal/content_loader"
	"github.com/kapitan123/telegrofler/bot/internal/content_loader/lovetik"
	"github.com/kapitan123/telegrofler/bot/internal/messenger"
	"github.com/kapitan123/telegrofler/bot/internal/storage"
	"github.com/kapitan123/telegrofler/bot/internal/systemclock"
	taskQueue "github.com/kapitan123/telegrofler/bot/internal/task_queue"
	"github.com/kapitan123/telegrofler/bot/internal/watermarker"

	"github.com/kapitan123/telegrofler/bot/config"
)

func main() {
	meta, err := config.NewMetadata()

	if err != nil {
		log.WithError(err).Fatal("Metada is not accessable")
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	client, err := firestore.NewClient(ctx, meta.GetProjectId())
	if err != nil {
		log.WithError(err).Fatal("Failed to create firestore client")
	}

	defer func() {
		err := client.Close()
		if err != nil {
			log.WithError(err).Fatal("Failed to close firestore client")
		}
	}()

	s := storage.New(client)

	botapi, err := tgbotapi.NewBotAPI(config.TelegramToken)
	if err != nil {
		log.WithError(err).Fatal("Failed to create bot api")
	}

	q := taskQueue.New(ctx, config.DeletionQueueName, meta)
	// AK TODO move bot and other shit to app
	// otherwise we need to create multiple instances of bot and storage to handle scheduler
	d := contentLoader.New(lovetik.New())
	m := messenger.New(botapi, d)
	w := watermarker.New()
	sc := systemclock.New()

	commandRunner := command.NewRunner(config.WorkersCount,
		command.WithRateLimit(choosePidor.New(m, s, w, q, sc)),
		command.WithRateLimit(toprofler.New(m, s)),
		command.WithRateLimit(toppidor.New(m, s)),
		recordBotPostReaction.New(m, s),
		recordReaction.New(m, s),
		replaceLinkWithMessage.New(m, s, d),
		replyTo300.New(m, q),
		replyToNo.New(m, w, q),
		replyToYes.New(m, q),
	)

	router := mux.NewRouter()

	// AK TODO pass args through app?
	setupRouter(router, commandRunner, choosePidor.New(m, s, w, q, sc), m)

	commandRunner.Start(ctx)

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", config.ServerPort),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.WithError(err).Fatal("Failed to start server")
	}
}
