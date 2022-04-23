package roflers

import (
	"time"
)

type Rofler struct {
	UserName  string     `firestore:"user_name"`
	Accolades []Accolade `firestore:"accolades"`
	UserPic   []byte     `firestore:"user_pic"`
}

type Accolade struct {
	Type      string    `firestore:"type"`       // rofler, worker, taste
	TimeFrame string    `firestore:"time_frame"` // week, month, year
	WonOn     time.Time `firestore:"won_on"`
}

// func (r *Rofler) AddPost(id, url string) {
// 	p := Post{
// 		VideoId:  id,
// 		Url:      url,
// 		PostedOn: time.Now(),
// 	}

// 	r.Posts = append(r.Posts, p)
// }
