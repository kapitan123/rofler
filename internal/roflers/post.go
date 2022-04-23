package roflers

import (
	"time"
)

// Video post stats for future important analytics
type Post struct {
	VideoId        string     `firestore:"video_id"`
	RoflerUserName string     `firestore:"rofler_user_name"`
	Url            string     `firestore:"url"`
	Reactions      []Reaction `firestore:"reactions"`
	KeyWords       []string   `firestore:"key_words"`
	PostedOn       time.Time  `firestore:"posted_on"`
}

type Reaction struct {
	Sender   string    `firestore:"sender"`
	Text     string    `firestore:"text"`
	PostedOn time.Time `firestore:"posted_on"`
}

func (p *Post) AddReaction(sender, text string) {
	reaction := Reaction{
		Sender:   sender,
		Text:     text,
		PostedOn: time.Now(),
	}

	p.Reactions = append(p.Reactions, reaction)
}
