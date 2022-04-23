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
	VideoId   string     `firestore:"video_id"`
	Url       string     `firestore:"url"`
	Reactions []Reaction `firestore:"reactions"`
	KeyWords  []string   `firestore:"key_words"`
	PostedOn  time.Time  `firestore:"posted_on"`
}

type Reaction struct {
	Sender string `firestore:"sender"`
	Text   string `firestore:"text"`
}

func (r *Rofler) AddPost(id string, url string) {
	p := Post{
		VideoId:  id,
		Url:      url,
		PostedOn: time.Now(),
	}

	r.Posts = append(r.Posts, p)
}

func (r *Rofler) AddReaction(vid string) {
	for _, p := range r.Posts {
		if p.VideoId == vid {
			p.ReactionsCount++
			//break
		}
	}

	// ak todo the problem is in mutbility of a struct
	_ = r.Posts
}

type RoflerStore interface {
	GetAll() ([]Rofler, error)
	Upsert(Rofler) error
	GetByUserName(string) (Rofler, error)
	//GetTop(string, string) (Rofler, error)
}
