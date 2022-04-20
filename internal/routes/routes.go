package routes

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/roflers"
	"github.com/kapitan123/telegrofler/pkg/problem"
	"github.com/kapitan123/telegrofler/pkg/tiktok"

	log "github.com/sirupsen/logrus"
)

// AK TODO should be in a separate aggreagate
type API struct {
	TikTok       *tiktok.TikTokClient
	Bot          *bot.Bot
	RoflersStore *roflers.RoflersStore
	// AK TODO add base concerns like liviness probe
}

func (api API) AddRoutes(router *mux.Router) {
	//router.HandleFunc("/stats", api.getAllStats).Methods("GET")
	//router.HandleFunc("/stats/top/{category}/{timespan}", api.callBackHandler).Methods("GET")
	router.HandleFunc("/callback", api.handleCallback).Methods("POST")
	//router.HandleFunc("/submit", api.callBackHandler).Methods("POST")
	router.HandleFunc("/download", api.download).Methods("GET")
}

// Handles callback from Telegram. Extracts url from message, converts video and uploads it back.
func (api API) handleCallback(resp http.ResponseWriter, req *http.Request) {
	wasHandeled, _ := api.replaceLinkWithMessage(req)
	if wasHandeled {
		api.captureReaction(req)
	}
}

// Downloads video from url and returns it as mp4 file compitable with telegram.
func (api API) download(resp http.ResponseWriter, req *http.Request) {
	tiktokUrl, err := url.QueryUnescape(req.URL.Query().Get("url"))

	if err != nil {
		prob := err.(*problem.Problem)
		prob.WriteTo(resp)
		return
	}

	log.Info("API: downloading video from ", tiktokUrl)

	item, err := api.TikTok.GetItemByUrl(tiktokUrl)

	// AK TODO should wrap it in a service
	if err != nil {
		prob := err.(*problem.Problem)
		prob.WriteTo(resp)
		return
	}

	content, err := api.TikTok.DownloadVideoFromItem(item)

	if err != nil {
		prob := err.(*problem.Problem)
		prob.WriteTo(resp)
		return
	}

	resp.Header().Set("Content-Type", "video/mp4")
	_, _ = resp.Write(content)
}

func (api API) replaceLinkWithMessage(req *http.Request) (bool, error) {
	wasHandeled := false
	tvp, err := api.Bot.TryExtractTikTokUrlData(req)

	if tvp == nil {
		return wasHandeled, nil
	}

	wasHandeled = true
	if err != nil {
		return wasHandeled, err
	}

	log.Info("Url was found in a callback message: ", tvp.Url)

	item, err := api.TikTok.GetItemByUrl(tvp.Url)

	if err != nil {
		return wasHandeled, err
	}

	bc, err := api.TikTok.DownloadVideoFromItem(item)

	if err != nil {
		return wasHandeled, err
	}

	tvp.VideoData.Payload = bc
	tvp.VideoData.Duration = item.Video.Duration
	tvp.VideoData.LikesCount = item.Stats.DiggCount
	tvp.VideoData.Title = item.Desc
	tvp.VideoData.Id = item.Id

	err = api.Bot.PostTiktokVideoFromUrl(tvp)

	if err != nil {
		return true, err
	}

	err = api.Bot.DeletePost(tvp.ChatId, tvp.OriginalMessageId)

	if err != nil {
		return wasHandeled, err
	}

	rofler, found, err := api.RoflersStore.GetByUserName(tvp.Sender)

	if err != nil {
		return wasHandeled, err
	}

	newPost := roflers.Post{
		TiktokId: tvp.VideoData.Id,
		Url:    :  tvp.Url,
		PostedOn :time.Now(),
	}

	rofler.AddPost(newPost)

	if !found {
		rofler.UserName = tvp.Sender
	}

	err = api.RoflersStore.Upsert(rofler)

	if err != nil {
		return wasHandeled, err
	}

	return wasHandeled, nil
}

func (api API) captureReaction(req *http.Request) error {
	tiktokid, sender, err := api.Bot.TryExtractTikTokReaction(req)

	if err != nil {
		return err
	}

	if tiktokid == "" {
		return nil
	}

	log.Infof("Reaction was found for %s sent by %s", tiktokid, sender)

	rofler, _, err := api.RoflersStore.GetByUserName(sender)

	if err != nil {
		return err
	}

	rofler.IncrementLike(tiktokid)

	err = api.RoflersStore.Upsert(rofler)

	if err != nil {
		return err
	}

	return nil
}
