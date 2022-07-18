package storage

import (
	"context"
	"time"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
)

// Video post stats for future important analytics
type (
	Post struct {
		VideoId        string     `firestore:"video_id"`
		Source         string     `firestore:"source"`
		RoflerUserName string     `firestore:"rofler_user_name"`
		RoflerId       int64      `firestore:"rofler_Id"`
		Url            string     `firestore:"url"`
		Reactions      []Reaction `firestore:"reactions"`
		KeyWords       []string   `firestore:"key_words"`
		PostedOn       time.Time  `firestore:"posted_on"`
		ChatId         int64      `firestore:"chat_id"`
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
	iter := s.client.Collection(postsCollection).Documents(ctx)
	return takeAll[Post](iter)
}

func (s *Storage) GetLastWeekPosts(ctx context.Context) ([]Post, error) {
	query := s.client.Collection(postsCollection).Where("posted_on", ">", time.Now().AddDate(0, 0, -7))

	iter := query.Documents(ctx)

	return takeAll[Post](iter)
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

func takeAll[T Post | Pidor](iter *firestore.DocumentIterator) ([]T, error) {
	var records []T
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var r T
		if err := doc.DataTo(&r); err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	return records, nil
}
