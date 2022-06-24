package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/command"
)

func setupRouter(r *mux.Router, runner *command.Runner) {
	r.HandleFunc("/callback", messageHandler(runner)).Methods("POST")
	// For scheduling
	//r.HandleFunc("/chat/rofler/top/{week}", app.getTopRoflerHandler).Methods("POST")
	//r.HandleFunc("/chat/gayoftheday", app.getTopRoflerHandler).Methods("POST")
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
