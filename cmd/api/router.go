package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/command"
)

// AK TODO we just can't pass all stuff here, we still need an abstraction to group configuration
func setupRouter(r *mux.Router, runner *command.Runner) {
	r.HandleFunc("/callback", messageHandler(runner)).Methods("POST")
	//r.HandleFunc("/chat/rofler/top/{week}", app.getTopRoflerHandler).Methods("POST")
	r.HandleFunc("/chat/gayoftheday", runChoosePidorCommand()).Methods("PUT")
}

// Intentionally swallows all exception so messages are not resend
// AK TODO send messages to a dead message quee
func messageHandler(runner *command.Runner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var update tgbotapi.Update
		err := json.NewDecoder(r.Body).Decode(&update)
		if err != nil {
			log.Error("Failed to decode the callback message.")
			w.WriteHeader(http.StatusOK)
			return
		}

		logContent(update.Message)

		err = runner.Run(r.Context(), update.Message)
		if err != nil {
			log.Error("Failed trying to invoke a command.")
		}
		w.WriteHeader(http.StatusOK)
	}
}

func choosePidorHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		logContent(update.Message)

		err = runner.Run(r.Context(), update.Message)
		if err != nil {
			log.Error("Failed trying to invoke a command.")
		}
		w.WriteHeader(http.StatusOK)
	}
}

func logContent(update *tgbotapi.Message) {
	ujs, _ := json.Marshal(update)
	log.Info("Callback content:", string(ujs))
}
