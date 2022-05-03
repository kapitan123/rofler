package routes

import (
	"encoding/json"
	"net/http"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorilla/mux"
	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/bot/tgaction"
	"github.com/kapitan123/telegrofler/internal/source/sourceFactory"
	log "github.com/sirupsen/logrus"
)

func (api *API) AddRoutes(router *mux.Router) {
	router.HandleFunc("/stats/posts", api.getAllPosts).Methods("GET")
	router.HandleFunc("/stats/top", api.getTopRoflerHandler).Methods("GET")
	router.HandleFunc("/download", api.download).Methods("GET")
	router.HandleFunc("/callback", api.handleCallback).Methods("POST")
}

func (api *API) AddHandlers() {
	api.handlers = tgaction.InitHandlers(api.Bot, api.PostsStore)
	api.commands = tgaction.InitCommands(api.Bot, api.PostsStore)
}

// Downloads video from url and returns it as mp4 file compitable with telegram.
func (api *API) download(resp http.ResponseWriter, req *http.Request) {
	vidUrl, err := url.QueryUnescape(req.URL.Query().Get("url"))

	if err != nil {
		writeTo(err, resp)
		return
	}

	log.Info("API: downloading video from ", vidUrl)
	extract, found := sourceFactory.TryGetExtractor(vidUrl)

	if !found {
		resp.WriteHeader(http.StatusNotFound)
		message := make(map[string]string)
		message["message"] = "No handler for this url"
		jsonResp, err := json.Marshal(message)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		resp.Write(jsonResp)
		return
	}

	lti, err := extract(vidUrl)

	// AK TODO should wrap it in a service
	if err != nil {
		writeTo(err, resp)
		return
	}

	resp.Header().Set("Content-Type", "video/mp4")
	_, _ = resp.Write(lti.Payload)
}

func (api *API) getTopRoflerHandler(resp http.ResponseWriter, req *http.Request) {
	tr, _, err := api.GetTopRoflerFromPosts()
	if err != nil {
		writeTo(err, resp)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(tr)
	_, _ = resp.Write(json)
}

func (api *API) getAllPosts(resp http.ResponseWriter, req *http.Request) {
	roflers, err := api.GetAllPosts()

	if err != nil {
		writeTo(err, resp)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(roflers)
	_, _ = resp.Write(json)
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

// AK TODO quick fix
func writeTo(err error, resp http.ResponseWriter) {
	log.Infof("API error %s", err)

	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(err)
}
