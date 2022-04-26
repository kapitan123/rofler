package tghandler

import (
	//"github.com/kapitan123/telegrofler/internal/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/bot"
)

type ThreeHundredCheck struct {
	*bot.Bot
}

func NewThreeHundredCheck(b *bot.Bot) *RecordReaction {
	return &RecordReaction{
		Bot: b,
	}
}

func (h *ThreeHundredCheck) Handle(m *tgbotapi.Message) (bool, error) {
	isHandeled := m.Text == "300"

	// if there is 300 ot trista
	// post a reaction with a sticker

	// if err != nil {
	// 	return !isHandeled, err
	// }

	// if reaction.Sender == "" {
	// 	return isHandeled, nil
	// }

	// log.Infof("Reaction was found for %s sent by %s", reaction.VideoId, reaction.Sender)

	// h.AddReactionToPost(reaction)

	return isHandeled, nil
}
