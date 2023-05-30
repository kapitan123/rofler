package message

import (
	"fmt"
	"net/url"
	"regexp"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/service/bot/domain"
	media "github.com/kapitan123/telegrofler/service/bot/domain/media_type"
	"github.com/samber/lo"
)

var mobileTiktokRegex = regexp.MustCompile(`https:\/\/[a-zA-Z]{2}\.tiktok\.com\/`)
var youtubeShortsRegex = regexp.MustCompile(`https:\/\/youtube.com\/shorts\/.*`)
var instagramReelRegex = regexp.MustCompile(`https:\/\/www.instagram.com\/reel\/.*`)

var supportedMasks = []*regexp.Regexp{
	mobileTiktokRegex,
	youtubeShortsRegex,
	instagramReelRegex,
}

type Message struct {
	Id      int
	message *tgbotapi.Message
	rtm     *tgbotapi.Message
}

func New(message *tgbotapi.Message) Message {
	return Message{
		Id:      message.MessageID,
		message: message,
		rtm:     message.ReplyToMessage,
	}
}

func (m Message) From() domain.UserRef {
	return domain.NewUserRef(m.message.From.ID, m.message.From.FirstName, m.message.From.LastName)
}

func (m Message) ChatId() int64 {
	return m.message.Chat.ID
}

func (m Message) HasMedia() bool {
	return m.message.Video != nil || m.message.Photo != nil
}

func (m Message) IsReplyToMedia() bool {
	return m.rtm.Video != nil || m.rtm.Photo != nil
}

func (m Message) IsSelfReply() bool {
	return m.message.From.ID == m.rtm.From.ID
}

func (m Message) MediaType() media.Type {
	if m.rtm.Video != nil {
		return media.Video
	} else if len(m.rtm.Photo) > 0 {
		return media.Image
	}

	return media.Type{}
}

func (m Message) MediaId() string {
	mediaId := ""
	if m.rtm.Video != nil {
		mediaId = m.rtm.Video.FileID
	} else if len(m.rtm.Photo) > 0 {
		mediaId = m.rtm.Photo[0].FileUniqueID
	}

	return mediaId
}

func (m Message) FindUrl() (*url.URL, bool, error) {
	if len(m.message.Entities) < 1 {
		return nil, false, fmt.Errorf("message contains no urls")
	}

	urlEnt, found := lo.Find(m.message.Entities, func(ent tgbotapi.MessageEntity) bool {
		return ent.Type == "url"
	})

	if !found {
		return nil, false, fmt.Errorf("message contains no urls")
	}

	url, err := urlEnt.ParseURL()

	if err != nil {
		return nil, false, fmt.Errorf("can't parse url")
	}

	return url, true, nil
}

func (m Message) HasDownloadableUrl() bool {
	url, found, err := m.FindUrl()

	if err != nil {
		return false
	}

	if !found {
		return false
	}

	for _, regex := range supportedMasks {
		if regex.MatchString(url.String()) {
			return true
		}
	}

	return false
}

func (m Message) AsReaction() domain.Reaction {
	reaction := domain.Reaction{
		Reactor:          m.From(),
		ReactToMessageId: m.rtm.MessageID,
		Text:             m.message.Text,
		PostedOn:         time.Now(),
	}

	return reaction
}
