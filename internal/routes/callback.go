package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kapitan123/telegrofler/internal/bot/tghandler"
)

func (api API) AddCallback(router *mux.Router) {
	router.HandleFunc("/callback", api.handleCallback).Methods("POST")
	api.handlers = tghandler.InitHandlers(api.Bot, api.PostsStore)
	api.commands = tghandler.InitCommands(api.Bot, api.PostsStore)
}

// AK TODO after extraction f handlers move back to routes
// Handles callback from Telegram. Extracts url from message, converts video and uploads it back.
func (api API) handleCallback(resp http.ResponseWriter, req *http.Request) {
	// AK TODO extract a Handler interface
	upd, err := api.GetUpdate(req)
	mess := upd.Message

	if err != nil {
		writeTo(err, resp)
		return
	}

	if mess == nil {
		return
	}

	// it invokes only handlers
	// does not invoke commands
	for _, h := range api.handlers {
		wasHandled, err := h.Handle(mess)

		if err != nil {
			writeTo(err, resp)
			return
		}

		if wasHandled {
			return
		}
	}

	command, err := api.GetCommand(mess)

	if err != nil {
		writeTo(err, resp)
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
		writeTo(err, resp)
		return
	}
}
