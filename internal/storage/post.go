package storage

import (
	"context"
	"google.golang.org/api/iterator"
	"time"
)

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
	MessageId int       `firestore:"message_id"` // ReplyToMessage.ID not the update.Message.ID
	Sender    string    `firestore:"sender"`
	Text      string    `firestore:"text"`
	PostedOn  time.Time `firestore:"posted_on"`
}

const postsCollection = "posts"

func (s *Storage) GetAll(ctx context.Context) ([]Post, error) {
	docs := s.client.Collection(postsCollection).Documents(ctx)
	var posts []Post
	for {
		doc, err := docs.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, err
		}
		var post Post
		if err := doc.DataTo(&post); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
