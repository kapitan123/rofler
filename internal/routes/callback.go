package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/bot/tgaction"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func (api *API) AddCallback(router *mux.Router) {
	router.HandleFunc("/callback", api.handleCallback).Methods("POST")
	api.handlers = tgaction.InitHandlers(api.Bot, api.PostsStore)
	api.commands = tgaction.InitCommands(api.Bot, api.PostsStore)
}

func (api *API) handleCallback(resp http.ResponseWriter, req *http.Request) {
	upd, err := api.GetUpdate(req)
	mess := upd.Message

	if err != nil {
		log.Error(err)
		return
	}

	if mess == nil {
		return
	}

	go tryHandleMessage(api, mess)
}

func tryHandleMessage(api *API, mess *tgbotapi.Message) {
	for _, h := range *api.handlers {
		wasHandled, err := h.Handle(mess)

		if err != nil {
			log.Error(err)
			return
		}

		if wasHandled {
			log.Infof("Handled by: %T", h)
			return
		}
	}

	command, err := bot.GetCommand(mess)

	if err != nil {
		log.Error(err)
		return
	}

	if command == "" {
		return
	}

	cmd, found := api.commands[command]

	if !found {
		return
	}

	_, err = cmd.Handle(mess)

	if err != nil {
		log.Error(err)
		return
	}

	log.Infof("Handled by: %T", cmd)
}
