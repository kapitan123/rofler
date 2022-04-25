package routes

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/data/model"
	"github.com/kapitan123/telegrofler/pkg/source/sourceFactory"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func (api API) AddCallback(router *mux.Router) {
	router.HandleFunc("/callback", api.handleCallback).Methods("POST")
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

func (api API) tryReplaceLinkWithMessage(mess *tgbotapi.Message) (bool, error) {
	isHandeled := true
	source, found := sourceFactory.TryGetSource(mess.Text)
	if !found {
		return !isHandeled, nil
	}

	svp := api.Bot.ConvertToSourceVideoPost(mess)

	if svp == nil {
		return !isHandeled, nil
	}

	log.Info("Url was found in a callback message: ", svp.Url)

	evi, err := source.ExtractVideoFromUrl(svp.Url)

	if err != nil {
		return isHandeled, err
	}

	svp.VideoData.Payload = evi.Payload
	svp.VideoData.Title = evi.Title
	svp.VideoData.Id = evi.Id

	log.Info("Trying to post to telegram: ", svp)

	err = api.Bot.RepostConvertedVideo(svp)

	if err != nil {
		return isHandeled, err
	}

	// we don't really care if if has failed and it makes integration tests a lot easier
	_ = api.Bot.DeletePost(svp.ChatId, svp.OriginalMessageId)

	newPost := model.Post{
		VideoId:        svp.VideoData.Id,
		Source:         evi.Type,
		RoflerUserName: svp.Sender,
		Url:            svp.Url,
		Reactions:      []model.Reaction{},
		KeyWords:       []string{},
		PostedOn:       time.Now(),
	}

	api.RoflersStore.UpsertPost(newPost)

	if err != nil {
		return isHandeled, err
	}

	return isHandeled, nil
}

func (api API) tryRecordReaction(m *tgbotapi.Message) (bool, error) {
	isHandeled := true
	reaction, err := api.Bot.TryExtractVideoRepostReaction(m)
	if err != nil {
		return !isHandeled, err
	}

	if reaction.Sender == "" {
		return isHandeled, nil
	}

	log.Infof("Reaction was found for %s sent by %s", reaction.VideoId, reaction.Sender)

	api.RoflersStore.AddReactionToPost(reaction)

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

	err = api.Bot.PostTopRofler(m.Chat.ID, tr, roflCount)
	if err != nil {
		return false, err
	}

	return true, nil
}
