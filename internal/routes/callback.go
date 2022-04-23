package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kapitan123/telegrofler/internal/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func (api API) AddCallback(router *mux.Router) {
	router.HandleFunc("/callback", api.handleCallback).Methods("POST")
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
		_, err = api.tryRecordReaction(mess.ReplyToMessage)
	}

	if !wasHandled && err == nil {
		_, err = api.tryExecCommand(mess)
	}

	if err != nil {
		writeTo(err, resp)
		return
	}
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

	if !found {
		rofler.UserName = tvp.Sender
	}

	rofler.AddPost(tvp.VideoData.Id, tvp.Url)

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
