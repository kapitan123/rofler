package infra

import (
	"context"
	"time"
)

func (s *Storage) GetLastWeekPosts(ctx context.Context, chatid int64) ([]Post, error) {
	query := s.client.Collection(postsCollection).Where("posted_on", ">=", time.Now().AddDate(0, 0, -7))
	query = query.Where("chat_id", "==", chatid)

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
