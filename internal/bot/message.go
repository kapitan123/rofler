package bot

import (
	"regexp"
	"strings"

	_ "embed"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5" // https://go-telegram-bot-api.dev/
)

// AK TODO add sucess parameter remove it from bot add it as an extemsion
func ExtractUserMediaReaction(upd *tgbotapi.Message) (ReplyToMediaPost, error) {
	rtm := upd.ReplyToMessage
	vr := ReplyToMediaPost{
		VideoId: rtm.Video.FileID,
		Reaction: Reaction{
			Sender:    upd.From.UserName,
			MessageId: rtm.MessageID,
			Text:      upd.Text,
		},
	}

	return vr, nil
}

func GetCommand(m *tgbotapi.Message) (string, error) {
	// AK TODO handles only one command
	if !m.IsCommand() {
		return "", nil
	}

	if strings.HasPrefix("top", m.Command()) {
		return "top", nil
	}

	return "", nil
}

// Handles incoming chat messages.
// Tries to extract a TikTok video url from the message if no url was found returns nil
// Handles only mobile format
func (b *Bot) ConvertToSourceVideoPost(m *tgbotapi.Message) *SourceVideoPost {
	return &SourceVideoPost{
		Sender:            m.From.UserName,
		ChatId:            m.Chat.ID,
		Url:               m.Text,
		OriginalMessageId: m.MessageID,
	}
}

const posterMaker = `🔥@(.*?)🔥`

// AK TODO add sucess parameter
func TryExtractVideoRepostReaction(upd *tgbotapi.Message) (ReplyToMediaPost, error) {
	vr := ReplyToMediaPost{}
	rtm := upd.ReplyToMessage

	if rtm == nil || rtm.From.UserName != "TelegroflBot" || rtm.Video == nil {
		return vr, nil
	}

	r := regexp.MustCompile(posterMaker)
	poster := r.FindStringSubmatch(rtm.Caption)[1]

	sender := upd.From.UserName
	// if the user reference his own post it is not a reaction. Can be moved outside the scope
	if sender == poster {
		return vr, nil
	}

	reply := ReplyToMediaPost{
		VideoId: rtm.Video.FileName,
		Reaction: Reaction{
			Sender:    sender,
			Text:      upd.Text,
			MessageId: upd.MessageID,
		},
	}

	return reply, nil
}
