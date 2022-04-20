package routes

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/roflers"
	"github.com/kapitan123/telegrofler/pkg/problem"
	"github.com/kapitan123/telegrofler/pkg/tiktok"

	log "github.com/sirupsen/logrus"
)

type API struct {
	TikTok *tiktok.TikTokClient
	Bot    *bot.Bot
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
	api.replaceLinkWithMessage(req)
	// AK TODO add error handling
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

func (api API) replaceLinkWithMessage(req *http.Request) error {
	log.Info("Telegram callback was called")

	tvp, err := api.Bot.TryExtractTikTokUrlData(req)

	if tvp == nil {
		return nil
	}

	if err != nil {
		return err
	}

	log.Info("Url was found in a callback message: ", tvp.Url)

	item, err := api.TikTok.GetItemByUrl(tvp.Url)

	if err != nil {
		return err
	}

	bc, err := api.TikTok.DownloadVideoFromItem(item)

	if err != nil {
		return err
	}

	tvp.VideoData.Payload = bc
	tvp.VideoData.Duration = item.Video.Duration
	tvp.VideoData.LikesCount = item.Stats.DiggCount
	tvp.VideoData.Title = item.Desc
	tvp.VideoData.Id = item.Id

	err = api.Bot.PostTiktokVideoFromUrl(tvp)

	if err != nil {
		return err
	}

	err = api.Bot.DeletePost(tvp.ChatId, tvp.OriginalMessageId)

	if err != nil {
		return err
	}

	store := roflers.New()

	roflerU, found, err := store.GetByUserName(tvp.Sender)

	if err != nil {
		return err
	}

	newPost := roflers.Post{
		TiktokId:      tvp.VideoData.Id,
		Url:           tvp.Url,
		ChatLikeCount: tvp.VideoData.LikesCount,
	}

	if !found {
		roflerU.UserName = tvp.Sender
		roflerU.AddPost(newPost)
	}

	err = store.Upsert(roflerU)

	if err != nil {
		return err
	}

	return nil
}
