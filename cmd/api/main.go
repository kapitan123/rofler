package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	"github.com/kapitan123/telegrofler/internal/messenger"
	"github.com/kapitan123/telegrofler/internal/messenger/formatter"
	"github.com/kapitan123/telegrofler/internal/services/watermarker"
	"github.com/kapitan123/telegrofler/internal/storage"
	log "github.com/sirupsen/logrus"

	"github.com/kapitan123/telegrofler/config"
)

func main() {
	flag.Parse()
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	client, err := firestore.NewClient(ctx, config.ProjectId)
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
	m := messenger.New(botapi)
	w := watermarker.New()
	f := formatter.New()

	commandRunner := command.NewRunner(config.WorkersCount,
		choosePidor.New(m, s, w),
		recordBotPostReaction.New(m, s),
		recordReaction.New(m, s),
		replaceLinkWithMessage.New(m, s),
		replyTo300.New(m),
		replyToNo.New(m, w),
		replyToYes.New(m),
		toprofler.New(m, s, f),
		toppidor.New(m, s, f),
	)

	log.WithField("addr", config.ServerPort).Info("Starting server on :%d", config.ServerPort)

	router := mux.NewRouter()
	// AK TODO pass args through app?
	setupRouter(router, commandRunner, choosePidor.New(m, s, w))

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", config.ServerPort),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.WithError(err).Fatal("Failed to start server")
	}
}
