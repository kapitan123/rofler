package main

import (
	"context"
	"fmt"
	"github.com/kapitan123/telegrofler/internal/contentLoader/webm"
	"net/http"

	"cloud.google.com/go/firestore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/kapitan123/telegrofler/internal/command"
	"github.com/kapitan123/telegrofler/internal/command/choosePidor"
	"github.com/kapitan123/telegrofler/internal/command/recordBotPostReaction"
	"github.com/kapitan123/telegrofler/internal/command/recordReaction"
	"github.com/kapitan123/telegrofler/internal/command/replaceLinkWithMessage"
	"github.com/kapitan123/telegrofler/internal/command/replyTo300"
	"github.com/kapitan123/telegrofler/internal/command/replyToNo"
	"github.com/kapitan123/telegrofler/internal/command/replyToYes"
	"github.com/kapitan123/telegrofler/internal/command/statsTrend"
	"github.com/kapitan123/telegrofler/internal/command/toppidor"
	"github.com/kapitan123/telegrofler/internal/command/toprofler"
	"github.com/kapitan123/telegrofler/internal/contentLoader"
	"github.com/kapitan123/telegrofler/internal/contentLoader/lovetik"
	"github.com/kapitan123/telegrofler/internal/contentLoader/mp4"
	"github.com/kapitan123/telegrofler/internal/contentLoader/shortsget"
	"github.com/kapitan123/telegrofler/internal/messenger"
	"github.com/kapitan123/telegrofler/internal/storage"
	"github.com/kapitan123/telegrofler/internal/systemclock"
	"github.com/kapitan123/telegrofler/internal/taskQueue"
	"github.com/kapitan123/telegrofler/internal/watermarker"

	"github.com/kapitan123/telegrofler/config"
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
	d := contentLoader.New(shortsget.New(), lovetik.New(), mp4.New(), webm.New())
	m := messenger.New(botapi, d)
	w := watermarker.New()
	sc := systemclock.New()

	commandRunner := command.NewRunner(config.WorkersCount,
		command.WithRateLimit(choosePidor.New(m, s, w, q, sc)),
		command.WithRateLimit(toprofler.New(m, s)),
		command.WithRateLimit(toppidor.New(m, s)),
		statsTrend.New(m, s),
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
