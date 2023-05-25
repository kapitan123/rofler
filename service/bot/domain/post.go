package domain

import (
	"time"
)

type (
	Post struct {
		VideoId   string
		Source    string
		Url       string
		Reactions []Reaction
		PostedOn  time.Time
		ChatId    int64
		UserRef   UserRef
	}

	Reaction struct {
		ReactToMessageId int
		Sender           UserRef
		Text             string
		PostedOn         time.Time
		ReactorUserRef   UserRef
	}
)

func (p *Post) AddReaction(reactor UserRef, text string, messageid int) {
	reaction := Reaction{
		ReactorUserRef:   reactor,
		Text:             text,
		ReactToMessageId: messageid,
		PostedOn:         time.Now(),
	}

	p.Reactions = append(p.Reactions, reaction)
}
