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
	"github.com/kapitan123/telegrofler/internal/command/recordBotPostReaction"
	"github.com/kapitan123/telegrofler/internal/command/recordReaction"
	"github.com/kapitan123/telegrofler/internal/command/replaceLinkWithMessage"
	"github.com/kapitan123/telegrofler/internal/command/replyTo300"
	"github.com/kapitan123/telegrofler/internal/command/replyToNo"
	"github.com/kapitan123/telegrofler/internal/command/replyToYes"
	"github.com/kapitan123/telegrofler/internal/command/toprofler"
	"github.com/kapitan123/telegrofler/internal/messenger"
	"github.com/kapitan123/telegrofler/internal/storage"
	log "github.com/sirupsen/logrus"

	"github.com/kapitan123/telegrofler/config"
)

const workers = 1

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

	m := messenger.New(botapi)
	commandRunner := command.NewRunner(config.WorkersCount,
		toprofler.New(m, s),
		replyToNo.New(m),
		replyTo300.New(m),
		replyToYes.New(m),
		recordBotPostReaction.New(m, s),
		recordReaction.New(m, s),
		replaceLinkWithMessage.New(m, s),
	)

	log.WithField("addr", config.ServerPort).Info("Starting server on :%d", config.ServerPort)

	router := mux.NewRouter()
	setupRouter(router, commandRunner)

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", config.ServerPort),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.WithError(err).Fatal("Failed to start server")
	}
}
