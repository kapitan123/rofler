package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"flag"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/command"
	"github.com/kapitan123/telegrofler/internal/command/rofler"
	"github.com/kapitan123/telegrofler/internal/messenger"
	"github.com/kapitan123/telegrofler/internal/storage"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var (
	telegramToken      = flag.String("token", "", "Telegram bot token")
	firestoreProjectId = flag.String("project", "", "Firestore project id")
	addr               = flag.String("addr", ":8080", "Address to listen on")
	workers            = flag.Int("workers", 1, "Number of workers")
)

func main() {
	flag.Parse()
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	client, err := firestore.NewClient(ctx, *firestoreProjectId)
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

	botapi, err := tgbotapi.NewBotAPI(*telegramToken)
	if err != nil {
		log.WithError(err).Fatal("Failed to create bot api")
	}
	m := messenger.New(botapi)
	commandRunner := command.NewRunner(*workers,
		rofler.New(m, s),
	)

	mux := http.NewServeMux()
	setupRouter(mux, commandRunner)
	log.WithField("addr", *addr).Info("Starting server")
	err = http.ListenAndServe(*addr, mux)
	if err != nil {
		log.WithError(err).Fatal("Failed to start server")
	}
}
