package storage

import (
	"context"
	"time"

	"google.golang.org/api/iterator"
)

// Video post stats for future important analytics
type (
	Post struct {
		VideoId        string     `firestore:"video_id"`
		Source         string     `firestore:"source"`
		RoflerUserName string     `firestore:"rofler_user_name"`
		Url            string     `firestore:"url"`
		Reactions      []Reaction `firestore:"reactions"`
		KeyWords       []string   `firestore:"key_words"`
		PostedOn       time.Time  `firestore:"posted_on"`
	}

	Reaction struct {
		MessageId int       `firestore:"message_id"` // RepllyToMessage.ID not the update.Message.ID
		Sender    string    `firestore:"sender"`
		Text      string    `firestore:"text"`
		PostedOn  time.Time `firestore:"posted_on"`
	}
)

func (p *Post) AddReaction(sender, text string, messageid int) {
	reaction := Reaction{
		Sender:    sender,
		Text:      text,
		MessageId: messageid,
		PostedOn:  time.Now(),
	}

	p.Reactions = append(p.Reactions, reaction)
}

const postsCollection = "posts"

func (s *Storage) GetAllPosts(ctx context.Context) ([]Post, error) {
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

	for _, p := range posts {
		if p.RoflerUserName == "" || p.RoflerUserName == "tester" {
			p.RoflerUserName = "unknown"
			s.client.Collection(postsCollection).Doc(p.VideoId).Delete(ctx)
		}
	}

	return posts, nil
}

func (s *Storage) UpsertPost(ctx context.Context, p Post) error {
	doc := s.client.Collection(postsCollection).Doc(p.VideoId)
	_, err := doc.Set(ctx, p)

	return err
}

func (s *Storage) GetPostById(ctx context.Context, videoId string) (Post, bool, error) {
	var p Post
	doc := s.client.Collection(postsCollection).Doc(videoId)
	snap, err := doc.Get(ctx)

	if err != nil {
		return p, false, nil
	}

	if err := snap.DataTo(&p); err != nil {
		return p, false, err
	}

	return p, true, nil
}

func (s *Storage) CreatePost(ctx context.Context, p Post) error {
	doc := s.client.Collection(postsCollection).Doc(p.VideoId)
	_, err := doc.Create(ctx, p)

	return err
}
