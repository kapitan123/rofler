package main

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/command"
	"net/http"
)

func setupRouter(r *http.ServeMux, runner *command.Runner) {
	r.Handle("/message", messageHandler(runner))
}

func messageHandler(runner *command.Runner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var update tgbotapi.Update
		err := json.NewDecoder(r.Body).Decode(&update)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = runner.Run(r.Context(), update.Message)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
}
