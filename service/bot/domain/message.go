package domain

import (
	"net/url"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type ChatId int64

func (m ChatId) String() string {
	return strconv.FormatInt(int64(m), 10)
}

type MessageId int

func (m MessageId) String() string {
	return strconv.Itoa(int(m))
}

type Message struct {
	Id             int
	message        *tgbotapi.Message
	rtm            *tgbotapi.Message
	ReplytoMessage ReplytoMessage
}

func NewMessage(message *tgbotapi.Message) Message {
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

func (m Message) From() UserRef {
	return NewUserRef(m.message.From.ID, m.message.From.FirstName, m.message.From.LastName)
}

func (m Message) MessageId() MessageId {
	return MessageId(m.message.MessageID)
}

func (m Message) ChatId() ChatId {
	return ChatId(m.message.Chat.ID)
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

func (m Message) MediaType() MediaType {
	if m.rtm.Video != nil {
		return Video
	} else if len(m.rtm.Photo) > 0 {
		return Image
	}

	return MediaType{}
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

func (m Message) containsOnlyUrlString() bool {
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
	if !m.containsOnlyUrlString() {
		return nil, errors.Errorf("message is not a url only message")
	}

	return getSupportedUrl(m.message.Text)
}

func (m Message) HasDownloadableUrl() bool {
	if !m.containsOnlyUrlString() {
		return false
	}

	return canBeDownloaded(m.message.Text)
}

func (m Message) AsReaction() Reaction {
	reaction := Reaction{
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
func (m ReplytoMessage) GetUserRef() (UserRef, error) {
	if len(m.rtm.CaptionEntities) == 0 || m.rtm.CaptionEntities[0].User == nil {
		return UserRef{}, errors.Errorf("message has no user reference")
	}
	user := m.rtm.CaptionEntities[0].User

	userRef := NewUserRef(user.ID, user.FirstName, user.LastName)

	return userRef, nil
}
