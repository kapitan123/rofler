package routes

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/roflers"
	"github.com/kapitan123/telegrofler/pkg/lovetik"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	router.HandleFunc("/callback", api.handleCallback).Methods("POST")
	router.HandleFunc("/download", api.download).Methods("GET")
}

// Handles callback from Telegram. Extracts url from message, converts video and uploads it back.
func (api API) handleCallback(resp http.ResponseWriter, req *http.Request) {
	upd, err := api.Bot.GetUpdate(req)
	mess := upd.Message

	if err != nil {
		writeTo(err, resp)
		return
	}

	if mess == nil {
		return
	}

	wasHandled, err := api.tryReplaceLinkWithMessage(mess)

	if !wasHandled && err == nil {
		_, err = api.tryRecordReaction(mess)
	}

	if !wasHandled && err == nil {
		_, err = api.tryExecCommand(mess)
	}

	if err != nil {
		writeTo(err, resp)
		return
	}
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

func (api API) tryReplaceLinkWithMessage(mess *tgbotapi.Message) (bool, error) {
	isHandeled := true

	tvp, err := api.Bot.ExtractTikTokVideoPost(mess)

	if tvp == nil {
		return !isHandeled, nil
	}

	if err != nil {
		return isHandeled, err
	}

	log.Info("Url was found in a callback message: ", tvp.Url)

	lti, err := api.LoveTik.DownloadVideoFromUrl(tvp.Url)

	if err != nil {
		return isHandeled, err
	}

	tvp.VideoData.Payload = lti.Payload
	tvp.VideoData.Title = lti.Title
	tvp.VideoData.Id = lti.Id

	log.Info("Trying to post to telegram: ", tvp)

	err = api.Bot.PostTiktokVideoFromUrl(tvp)

	if err != nil {
		return isHandeled, err
	}

	// we don't really care if if has failed and it makes integration tests a lot easier
	_ = api.Bot.DeletePost(tvp.ChatId, tvp.OriginalMessageId)

	rofler, found, err := api.RoflersStore.GetByUserName(tvp.Sender)

	if err != nil {
		return isHandeled, err
	}

	newPost := roflers.Post{
		VideoId:  tvp.VideoData.Id,
		Url:      tvp.Url,
		PostedOn: time.Now(),
	}

	rofler.AddPost(newPost)

	if !found {
		rofler.UserName = tvp.Sender
	}

	err = api.RoflersStore.Upsert(rofler)

	if err != nil {
		return isHandeled, err
	}

	return isHandeled, nil
}

func (api API) tryRecordReaction(m *tgbotapi.Message) (bool, error) {
	isHandeled := true
	reaction, err := api.Bot.TryExtractTikTokReaction(m)
	if err != nil {
		return !isHandeled, err
	}

	if reaction.Sender == "" {
		return isHandeled, nil
	}

	log.Infof("Reaction was found for %s sent by %s", reaction.VideoId, reaction.Sender)

	api.RoflersStore.IncrementLike(reaction)

	return isHandeled, nil
}

func (api API) tryExecCommand(m *tgbotapi.Message) (bool, error) {
	command, err := api.Bot.GetCommand(m)
	if err != nil {
		return false, err
	}

	if command != bot.TopCommand {
		return true, nil
	}

	tr, roflCount, err := api.RoflersStore.GetTopRofler()
	if err != nil {
		return false, err
	}

	err = api.Bot.PostTopRofler(m.Chat.ID, tr.UserName, roflCount)
	if err != nil {
		return false, err
	}

	return true, nil
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
