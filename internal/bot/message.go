package bot

import (
	"regexp"
	"strings"

	"github.com/kapitan123/telegrofler/internal/data/post"

	_ "embed"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5" // https://go-telegram-bot-api.dev/
)

// AK TODO add sucess parameter remove it from bot add it as an extemsion
func ExtractUserMediaReaction(upd *tgbotapi.Message) (post.VideoReaction, error) {
	rtm := upd.ReplyToMessage
	vr := post.VideoReaction{}

	vr.Sender = upd.From.UserName
	vr.MessageId = rtm.MessageID
	vr.Text = upd.Text
	vr.VideoId = rtm.Video.FileID

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
func (b *Bot) TryExtractVideoRepostReaction(upd *tgbotapi.Message) (post.VideoReaction, error) {
	vr := post.VideoReaction{}
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

	return post.VideoReaction{Sender: sender, VideoId: rtm.Video.FileName, Text: upd.Text, MessageId: upd.MessageID}, nil
}
