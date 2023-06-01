package message

import (
	"net/url"
	"regexp"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/service/bot/domain"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

var mobileTiktokRegex = regexp.MustCompile(`https:\/\/[a-zA-Z]{2}\.tiktok\.com\/`)
var youtubeShortsRegex = regexp.MustCompile(`https:\/\/youtube\.com\/shorts\/.*`)
var instagramReelRegex = regexp.MustCompile(`https:\/\/www.instagram.com\/reel\/.*`)

var supportedMasks = []*regexp.Regexp{
	mobileTiktokRegex,
	youtubeShortsRegex,
	instagramReelRegex,
}

type Message struct {
	Id             int
	message        *tgbotapi.Message
	rtm            *tgbotapi.Message
	ReplytoMessage ReplytoMessage
}

func New(message *tgbotapi.Message) Message {
	if message == nil {
		message = &tgbotapi.Message{}
	}

	tgReplyTo := message.ReplyToMessage

	if tgReplyTo == nil {
		tgReplyTo = &tgbotapi.Message{}
	}

	return Message{
		Id:             message.MessageID,
		message:        message,
		rtm:            tgReplyTo,
		ReplytoMessage: ReplytoMessage{tgReplyTo},
	}
}

func (m Message) IsCommand(commandName string) bool {
	return m.message.IsCommand() && m.message.Command() == commandName
}

func (m Message) From() domain.UserRef {
	return domain.NewUserRef(m.message.From.ID, m.message.From.FirstName, m.message.From.LastName)
}

func (m Message) ChatId() int64 {
	return m.message.Chat.ID
}

func (m Message) IsBotPost() bool {
	return m.message.From.IsBot
}

func (m Message) HasMedia() bool {
	return m.message.Video != nil || m.message.Photo != nil
}

func (m Message) IsSelfReply() bool {
	return m.message.From.ID == m.rtm.From.ID
}

func (m Message) MediaType() domain.MediaType {
	if m.rtm.Video != nil {
		return domain.Video
	} else if len(m.rtm.Photo) > 0 {
		return domain.Image
	}

	return domain.MediaType{}
}

func (m Message) MediaId() string {
	mediaId := ""
	if m.rtm.Video != nil {
		mediaId = m.rtm.Video.FileName
	} else if len(m.rtm.Photo) > 0 {
		mediaId = m.rtm.Photo[0].FileUniqueID
	}

	return mediaId
}

func (m Message) ContainsOnlyUrlString() bool {
	if len(m.message.Entities) != 1 {
		return false
	}

	urlEnt, found := lo.Find(m.message.Entities, func(ent tgbotapi.MessageEntity) bool {
		return ent.Type == "url"
	})

	if !found {
		return false
	}

	if urlEnt.Length != len(m.message.Text) {
		return false
	}

	return true
}

func (m Message) GetEmbeddedUrl() (*url.URL, error) {
	if !m.ContainsOnlyUrlString() {
		return nil, errors.Errorf("message is not a url only message")
	}

	return url.Parse(m.message.Text)
}

func (m Message) HasDownloadableUrl() bool {
	found := m.ContainsOnlyUrlString()

	if !found {
		return false
	}

	url := m.message.Text
	for _, regex := range supportedMasks {
		if regex.MatchString(url) {
			return true
		}
	}

	return false
}

func (m Message) AsReaction() domain.Reaction {
	reaction := domain.Reaction{
		Reactor:  m.From(),
		Text:     m.message.Text,
		PostedOn: time.Now(),
	}

	return reaction
}

type ReplytoMessage struct {
	rtm *tgbotapi.Message
}

func (m ReplytoMessage) HasMedia() bool {
	return m.rtm.Video != nil || m.rtm.Photo != nil
}

func (m ReplytoMessage) IsPostedByBot() bool {
	return m.rtm.From.IsBot && m.rtm.From.UserName == "TelegroflBot"
}

// Based on the fact that bot posts always contain exactly one mention
func (m ReplytoMessage) GetUserRef() (domain.UserRef, error) {
	if len(m.rtm.CaptionEntities) == 0 || m.rtm.CaptionEntities[0].User == nil {
		return domain.UserRef{}, errors.Errorf("message has no user reference")
	}
	user := m.rtm.CaptionEntities[0].User

	userRef := domain.NewUserRef(user.ID, user.FirstName, user.LastName)

	return userRef, nil
}
