package routes

import (
	"context"
	"encoding/json"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorilla/mux"
	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/bot/tgaction"
	log "github.com/sirupsen/logrus"
)

func (app *App) AddRoutes(router *mux.Router) {
	router.HandleFunc("/stats/top", app.getTopRoflerHandler).Methods("GET")
	router.HandleFunc("/callback", app.handleCallback).Methods("POST")
}

func (app *App) AddHandlers() {
	app.handlers = tgaction.InitHandlers(app.Bot, app.PostsStorage)
	app.commands = tgaction.InitCommands(app.Bot, app.PostsStorage)
}

func (app *App) getTopRoflerHandler(resp http.ResponseWriter, req *http.Request) {
	tr, _, err := app.PostsStorage.GetTopRoflerFromPosts(req.Context())
	if err != nil {
		writeTo(err, resp)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(tr)
	_, _ = resp.Write(json)
}

func (app *App) handleCallback(resp http.ResponseWriter, req *http.Request) {
	upd, err := app.GetUpdate(req)
	mess := upd.Message

	if err != nil {
		log.Error(err)
		return
	}

	if mess == nil {
		return
	}

	go tryHandleMessage(app, mess, req.Context())
}

func tryHandleMessage(app *App, mess *tgbotapi.Message, ctx context.Context) {
	for _, h := range *app.handlers {
		wasHandled, err := h.Handle(mess, ctx)

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

	cmd, found := app.commands[command]

	if !found {
		return
	}

	_, err = cmd.Handle(mess, ctx)

	if err != nil {
		log.Error(err)
		return
	}

	log.Infof("Handled by: %T", cmd)
}

// AK TODO quick fix
func writeTo(err error, resp http.ResponseWriter) {
	log.Infof("API error %s", err)

	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(err)
}
