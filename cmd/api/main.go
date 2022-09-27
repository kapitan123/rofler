package main

import (
	"context"
	"fmt"
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
	"github.com/kapitan123/telegrofler/internal/command/toppidor"
	"github.com/kapitan123/telegrofler/internal/command/toprofler"
	"github.com/kapitan123/telegrofler/internal/contentLoader"
	"github.com/kapitan123/telegrofler/internal/contentLoader/lovetik"
	"github.com/kapitan123/telegrofler/internal/contentLoader/mp4"
	"github.com/kapitan123/telegrofler/internal/contentLoader/shortsget"
	"github.com/kapitan123/telegrofler/internal/messenger"
	"github.com/kapitan123/telegrofler/internal/storage"
	"github.com/kapitan123/telegrofler/internal/systemclock"
	"github.com/kapitan123/telegrofler/internal/watermarker"

	"github.com/kapitan123/telegrofler/config"
)

func main() {
	meta, err := config.GetMetadata()

	if err != nil {
		log.WithError(err).Fatal("Metada is not accessable")
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	client, err := firestore.NewClient(ctx, meta.ProjectId)
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

	// AK TODO move bot and other shit to app
	// otherwise we need to create multiple instances of bot and storage to handle scheduler
	d := contentLoader.New(shortsget.New(), lovetik.New(), mp4.New())
	m := messenger.New(botapi, d)
	w := watermarker.New()
	sc := systemclock.New()

	commandRunner := command.NewRunner(config.WorkersCount,
		command.WithRateLimit(choosePidor.New(m, s, w, sc)),
		command.WithRateLimit(toprofler.New(m, s)),
		command.WithRateLimit(toppidor.New(m, s)),
		recordBotPostReaction.New(m, s),
		recordReaction.New(m, s),
		replaceLinkWithMessage.New(m, s, d),
		replyTo300.New(m),
		replyToNo.New(m, w),
		replyToYes.New(m),
	)

	router := mux.NewRouter()
	// AK TODO pass args through app?
	setupRouter(router, commandRunner, choosePidor.New(m, s, w, sc))

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
