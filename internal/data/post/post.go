package post

import (
	"time"
)

// Video post stats for future important analytics
type Post struct {
	VideoId        string     `firestore:"video_id"`
	Source         string     `firestore:"source"`
	RoflerUserName string     `firestore:"rofler_user_name"`
	Url            string     `firestore:"url"`
	Reactions      []Reaction `firestore:"reactions"`
	KeyWords       []string   `firestore:"key_words"`
	PostedOn       time.Time  `firestore:"posted_on"`
}

type Reaction struct {
	MessageId int       `firestore:"message_id"`
	Sender    string    `firestore:"sender"`
	Text      string    `firestore:"text"`
	PostedOn  time.Time `firestore:"posted_on"`
}

func (p *Post) AddReaction(sender, text string, messageid int) {
	reaction := Reaction{
		Sender:    sender,
		Text:      text,
		MessageId: messageid,
		PostedOn:  time.Now(),
	}

	p.Reactions = append(p.Reactions, reaction)
}
