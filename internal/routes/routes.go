package routes

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/bot/tgaction"
	"github.com/kapitan123/telegrofler/internal/data/post"
	"github.com/kapitan123/telegrofler/pkg/source/sourceFactory"

	log "github.com/sirupsen/logrus"
)

// AK TODO should be in a separate file
type API struct {
	*bot.Bot
	*post.PostsStore
	handlers []tgaction.BotMessageHandler
	commands map[string]tgaction.BotCommandHandler
}

func (api API) AddRoutes(router *mux.Router) {
	router.HandleFunc("/stats/posts", api.getAllPosts).Methods("GET")
	router.HandleFunc("/stats/top", api.getTopRoflerHandler).Methods("GET")
	router.HandleFunc("/download", api.download).Methods("GET")
}

// Downloads video from url and returns it as mp4 file compitable with telegram.
func (api API) download(resp http.ResponseWriter, req *http.Request) {
	vidUrl, err := url.QueryUnescape(req.URL.Query().Get("url"))

	if err != nil {
		writeTo(err, resp)
		return
	}

	log.Info("API: downloading video from ", vidUrl)
	source, found := sourceFactory.TryGetSource(vidUrl)

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

	lti, err := source.ExtractVideoFromUrl(vidUrl)

	// AK TODO should wrap it in a service
	if err != nil {
		writeTo(err, resp)
		return
	}

	resp.Header().Set("Content-Type", "video/mp4")
	_, _ = resp.Write(lti.Payload)
}

func (api API) getTopRoflerHandler(resp http.ResponseWriter, req *http.Request) {
	tr, _, err := api.GetTopRoflerFromPosts()
	if err != nil {
		writeTo(err, resp)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(tr)
	_, _ = resp.Write(json)
}

func (api API) getAllPosts(resp http.ResponseWriter, req *http.Request) {
	roflers, err := api.GetAllPosts()

	if err != nil {
		writeTo(err, resp)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(roflers)
	_, _ = resp.Write(json)
}

// AK TODO quick fix
func writeTo(err error, resp http.ResponseWriter) {
	log.Infof("API error %s", err)

	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(err)
}
