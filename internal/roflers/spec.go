package roflers

import (
	"time"
)

type Rofler struct {
	UserName string `firestore:"user_name"`
	Posts    []Post `firestore:"posts"`
}

// Video post stats for future important analytics
type Post struct {
	VideoId        string    `firestore:"video_id"`
	Url            string    `firestore:"url"`
	ReactionsCount int       `firestore:"reactions_count"`
	KeyWords       []string  `firestore:"key_words"`
	PostedOn       time.Time `firestore:"posted_on"`
}

func (r *Rofler) AddPost(p Post) {
	r.Posts = append(r.Posts, p)
}

func (r *Rofler) AddReaction(vid string) {
	for _, p := range r.Posts {
		if p.VideoId == vid {
			p.ReactionsCount++
		}
	}
}

type RoflerStore interface {
	GetAll() ([]Rofler, error)
	Upsert(Rofler) error
	GetByUserName(string) (Rofler, error)
	//GetTop(string, string) (Rofler, error)
}
