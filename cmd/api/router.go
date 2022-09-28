package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/command"
	"github.com/kapitan123/telegrofler/internal/command/choosePidor"
	"github.com/kapitan123/telegrofler/internal/messenger"
)

// AK TODO we just can't pass all stuff here, we still need an abstraction to group configuration
// temp solution with direct handler function
func setupRouter(r *mux.Router, runner *command.Runner, pdr *choosePidor.ChoosePidor, msgr *messenger.Messenger) {
	r.HandleFunc("/callback", messageHandler(runner)).Methods("POST")
	r.HandleFunc("/chat/{chatid}/pidoroftheday", choosePidorHandler(pdr)).Methods("POST")
	r.HandleFunc("/chat/{chatid}/{messageid}", deleteMessageHandler(msgr)).Methods("DELETE")
}

// AK TODO send messages to a dead message quee
func messageHandler(runner *command.Runner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var update tgbotapi.Update
		err := json.NewDecoder(r.Body).Decode(&update)
		if err != nil {
			log.WithError(err).Error("Failed to decode the callback message")

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		logContent(&update)

		err = runner.Run(r.Context(), update.Message)
		if err != nil {
			log.WithError(err).Error("Failed trying to invoke a command")
		}
		// Intentionally swallows all exception so messages are not resend
		w.WriteHeader(http.StatusOK)
	}
}

func choosePidorHandler(pdr *choosePidor.ChoosePidor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatId, err := extractChatId(r)
		if err != nil {
			log.WithError(err).Error("Failed trying to parse chat id", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		err = pdr.ChoosePidor(r.Context(), chatId)

		if err != nil {
			log.Error("Failed trying to choose a pidor", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	}
}

func deleteMessageHandler(msgr *messenger.Messenger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatId, err := extractChatId(r)
		if err != nil {
			log.WithError(err).Error("Failed trying to parse chat id", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		messageid, err := extractMessageId(r)
		if err != nil {
			log.WithError(err).Error("Failed trying to parse message id", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		err = msgr.Delete(chatId, messageid)

		if err != nil {
			log.Error("Failed trying to delete a message", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	}
}

func logContent(update *tgbotapi.Update) {
	ujs, _ := json.Marshal(update)
	log.Info("Callback content:", string(ujs))
}

func extractChatId(r *http.Request) (int64, error) {
	chatarg := mux.Vars(r)["chatid"]

	return strconv.ParseInt(chatarg, 10, 64)
}

func extractMessageId(r *http.Request) (int, error) {
	messagearg := mux.Vars(r)["messageid"]

	return strconv.Atoi(messagearg)
}
