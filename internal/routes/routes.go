package routes

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/roflers"
	"github.com/kapitan123/telegrofler/pkg/lovetik"

	log "github.com/sirupsen/logrus"
)

// AK TODO should be in a separate aggreagate
type API struct {
	//TikTok       *tiktok.TikTokClient
	LoveTik      *lovetik.LoveTikClient
	Bot          *bot.Bot
	RoflersStore *roflers.RoflersStore
	// AK TODO add base concerns like liviness probe
}

func (api API) AddRoutes(router *mux.Router) {
	router.HandleFunc("/stats", api.getAllStats).Methods("GET")
	router.HandleFunc("/stats/top", api.getTopRoflerHandler).Methods("GET")
	router.HandleFunc("/download", api.download).Methods("GET")
}

// Downloads video from url and returns it as mp4 file compitable with telegram.
func (api API) download(resp http.ResponseWriter, req *http.Request) {
	tiktokUrl, err := url.QueryUnescape(req.URL.Query().Get("url"))

	if err != nil {
		writeTo(err, resp)
		return
	}

	log.Info("API: downloading video from ", tiktokUrl)

	lti, err := api.LoveTik.DownloadVideoFromUrl(tiktokUrl)

	// AK TODO should wrap it in a service
	if err != nil {
		writeTo(err, resp)
		return
	}

	resp.Header().Set("Content-Type", "video/mp4")
	_, _ = resp.Write(lti.Payload)
}

func (api API) getTopRoflerHandler(resp http.ResponseWriter, req *http.Request) {
	tr, _, err := api.RoflersStore.GetTopRofler()
	if err != nil {
		writeTo(err, resp)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(tr)
	_, _ = resp.Write(json)
}

func (api API) getAllStats(resp http.ResponseWriter, req *http.Request) {
	roflers, err := api.RoflersStore.GetAll()

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
